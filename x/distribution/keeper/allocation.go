package keeper

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AllocateTokens performs reward and fee distribution to all validators based
// on the F1 fee distribution specification.
func (k Keeper) AllocateTokens(ctx context.Context, totalPreviousPower int64, bondedVotes []abci.VoteInfo) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	logger := k.Logger(sdkCtx)
	// fetch and clear the collected fees for distribution, since this is
	// called in BeginBlock, collected fees will be from the previous block
	// (and distributed to the previous proposer)
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// transfer collected fees to the distribution module account
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, feesCollectedInt)
	if err != nil {
		return err
	}

	if len(feesCollectedInt) > 0 {
		ratio, err := k.GetRatio(sdk.UnwrapSDKContext(ctx))
		if err != nil {
			return err
		}
		remaining := feesCollected

		// community tax
		communityTax, err := k.GetCommunityTax(ctx)
		if err != nil {
			return err
		}
		voteMultiplier := math.LegacyOneDec().Sub(communityTax)
		feeMultiplier := feesCollected.MulDecTruncate(voteMultiplier)

		// burn fee: ratio.Burn
		burnFee := k.CalculatePercentage(feeMultiplier, ratio.Burn)
		burnFeeDec := sdk.NewDecCoinsFromCoins(burnFee...)
		err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnFee)
		if err != nil {
			panic(err)
		}
		remaining = remaining.Sub(burnFeeDec)
		err = k.AddTotalBurned(ctx, burnFee)
		if err != nil {
			return err
		}

		// emit burn fee
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeBurnFee,
				sdk.NewAttribute(sdk.AttributeKeyAmount, burnFee.String()),
			),
		)

		// base fee: ratio.Base
		var baseFee sdk.Coins
		var baseFeeDec sdk.DecCoins
		base, err := k.GetBaseAddress(sdkCtx)
		if err == nil && base.Address != "" {
			baseAddr, err := sdk.AccAddressFromBech32(base.Address)
			if err != nil {
				return err
			}
			baseFee = k.CalculatePercentage(feeMultiplier, ratio.Base)
			baseFeeDec = sdk.NewDecCoinsFromCoins(baseFee...)
			err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, baseAddr, baseFee)
			if err != nil {
				panic(err)
			}
			remaining = remaining.Sub(baseFeeDec)
			// emit base fee
			sdkCtx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeBaseFee,
					sdk.NewAttribute(sdk.AttributeKeyAmount, baseFee.String()),
				),
			)
		} else {
			logger.Info("No base address found, base proportion will be added to staking rewards")
			ratio.StakingRewards = ratio.StakingRewards.Add(ratio.Base)
		}

		// feesCollectedInt = feesCollectedInt.Sub(burnFee...).Sub(baseFee...)
		feeMultiplier = feeMultiplier.Sub(burnFeeDec).Sub(baseFeeDec)

		// temporary workaround to keep CanWithdrawInvariant happy
		// general discussions here: https://github.com/cosmos/cosmos-sdk/issues/2906#issuecomment-441867634
		feePool, err := k.FeePool.Get(ctx)
		if err != nil {
			return err
		}

		if totalPreviousPower == 0 {
			feePool.CommunityPool = feePool.CommunityPool.Add(feesCollected...)
			return k.FeePool.Set(ctx, feePool)
		}

		// allocate tokens proportionally to voting power
		//
		// TODO: Consider parallelizing later
		//
		// Ref: https://github.com/cosmos/cosmos-sdk/pull/3099#discussion_r246276376
		for _, vote := range bondedVotes {
			validator, err := k.stakingKeeper.ValidatorByConsAddr(ctx, vote.Validator.Address)
			if err != nil {
				return err
			}

			// TODO: Consider micro-slashing for missing votes.
			//
			// Ref: https://github.com/cosmos/cosmos-sdk/issues/2525#issuecomment-430838701
			powerFraction := math.LegacyNewDec(vote.Validator.Power).QuoTruncate(math.LegacyNewDec(totalPreviousPower))
			reward := feeMultiplier.MulDecTruncate(powerFraction)

			err = k.AllocateTokensToValidator(ctx, validator, reward)
			if err != nil {
				return err
			}

			remaining = remaining.Sub(reward)
		}

		// allocate community funding
		feePool.CommunityPool = feePool.CommunityPool.Add(remaining...)

		// log
		logger.Info("Fee Distribution",
			types.EventTypeTotalFees, feesCollected.String(),
			types.EventTypeBurnFee, burnFeeDec.String(),
			types.EventTypeBaseFee, baseFeeDec.String(),
			types.EventTypeStakingRewards, feeMultiplier.String(),
			types.EventTypeCommunityPool, remaining.String())

		return k.FeePool.Set(ctx, feePool)
	}
	return nil
}

// AllocateTokensToValidator allocate tokens to a particular validator,
// splitting according to commission.
func (k Keeper) AllocateTokensToValidator(ctx context.Context, val stakingtypes.ValidatorI, tokens sdk.DecCoins) error {
	// split tokens between validator and delegators according to commission
	commission := tokens.MulDec(val.GetCommission())
	shared := tokens.Sub(commission)

	valBz, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
	if err != nil {
		return err
	}

	// update current commission
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator()),
		),
	)
	currentCommission, err := k.GetValidatorAccumulatedCommission(ctx, valBz)
	if err != nil {
		return err
	}

	currentCommission.Commission = currentCommission.Commission.Add(commission...)
	err = k.SetValidatorAccumulatedCommission(ctx, valBz, currentCommission)
	if err != nil {
		return err
	}

	// update current rewards
	currentRewards, err := k.GetValidatorCurrentRewards(ctx, valBz)
	if err != nil {
		return err
	}

	currentRewards.Rewards = currentRewards.Rewards.Add(shared...)
	err = k.SetValidatorCurrentRewards(ctx, valBz, currentRewards)
	if err != nil {
		return err
	}

	// update outstanding rewards
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, tokens.String()),
			sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator()),
		),
	)

	outstanding, err := k.GetValidatorOutstandingRewards(ctx, valBz)
	if err != nil {
		return err
	}

	outstanding.Rewards = outstanding.Rewards.Add(tokens...)
	return k.SetValidatorOutstandingRewards(ctx, valBz, outstanding)
}

// sendCommunityPoolToExternalPool does the following:
//
//	truncate the community pool value (DecCoins) to sdk.Coins
//	distribute from the distribution module account to the external community pool account
//	update the bookkept value in x/distribution
func (k Keeper) sendCommunityPoolToExternalPool(ctx sdk.Context) error {
	feePool, err := k.FeePool.Get(ctx)
	if err != nil {
		return err
	}

	if feePool.CommunityPool.IsZero() {
		return nil
	}

	amt, remaining := feePool.CommunityPool.TruncateDecimal()
	ctx.Logger().Debug(
		"sending distribution community pool amount to external pool pool",
		"pool", k.externalCommunityPool.GetCommunityPoolModule(),
		"amount", amt.String(),
		"remaining", remaining.String(),
	)
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.externalCommunityPool.GetCommunityPoolModule(), amt); err != nil {
		return err
	}

	return k.FeePool.Set(ctx, types.FeePool{CommunityPool: remaining})
}

func (k Keeper) CalculatePercentage(coins sdk.DecCoins, percentage math.LegacyDec) sdk.Coins {
	var result sdk.Coins
	for _, coin := range coins {
		// Multiply decimal amount by percentage, then truncate to integer
		amount := coin.Amount.MulTruncate(percentage).TruncateInt()
		if amount.IsPositive() {
			result = result.Add(sdk.NewCoin(coin.Denom, amount))
		}
	}
	return result
}
