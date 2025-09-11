package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Verify interface at compile time
var (
	_ sdk.Msg = (*MsgSetWithdrawAddress)(nil)
	_ sdk.Msg = (*MsgWithdrawDelegatorReward)(nil)
	_ sdk.Msg = (*MsgWithdrawValidatorCommission)(nil)
	_ sdk.Msg = (*MsgUpdateParams)(nil)
	_ sdk.Msg = (*MsgCommunityPoolSpend)(nil)
	_ sdk.Msg = (*MsgDepositValidatorRewardsPool)(nil)
	_ sdk.Msg = (*MsgChangeRatio)(nil)
	_ sdk.Msg = (*MsgChangeBaseAddress)(nil)
	_ sdk.Msg = (*MsgChangeModerator)(nil)
	_ sdk.Msg = (*MsgResetTotalBurned)(nil)
)

func NewMsgSetWithdrawAddress(delAddr, withdrawAddr sdk.AccAddress) *MsgSetWithdrawAddress {
	return &MsgSetWithdrawAddress{
		DelegatorAddress: delAddr.String(),
		WithdrawAddress:  withdrawAddr.String(),
	}
}

func NewMsgWithdrawDelegatorReward(delAddr, valAddr string) *MsgWithdrawDelegatorReward {
	return &MsgWithdrawDelegatorReward{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
	}
}

func NewMsgWithdrawValidatorCommission(valAddr string) *MsgWithdrawValidatorCommission {
	return &MsgWithdrawValidatorCommission{
		ValidatorAddress: valAddr,
	}
}

// NewMsgFundCommunityPool returns a new MsgFundCommunityPool with a sender and
// a funding amount.
func NewMsgFundCommunityPool(amount sdk.Coins, depositor string) *MsgFundCommunityPool {
	return &MsgFundCommunityPool{
		Amount:    amount,
		Depositor: depositor,
	}
}

// NewMsgDepositValidatorRewardsPool returns a new MsgDepositValidatorRewardsPool
// with a depositor and a funding amount.
func NewMsgDepositValidatorRewardsPool(depositor, valAddr string, amount sdk.Coins) *MsgDepositValidatorRewardsPool {
	return &MsgDepositValidatorRewardsPool{
		Amount:           amount,
		Depositor:        depositor,
		ValidatorAddress: valAddr,
	}
}

// NewMsgChangeRatio returns a new MsgChangeRatio with a new distribution ratio
func NewMsgChangeRatio(moderator sdk.AccAddress, ratio Ratio) *MsgChangeRatio {
	return &MsgChangeRatio{
		ModeratorAddress: moderator.String(),
		Ratio:            ratio,
	}
}

// NewMsgChangeBaseAddress returns a new MsgChangeBaseAddress with a new base address
func NewMsgChangeBaseAddress(moderator sdk.AccAddress, newBaseAddress sdk.AccAddress) *MsgChangeBaseAddress {
	return &MsgChangeBaseAddress{
		ModeratorAddress: moderator.String(),
		NewBaseAddress:   newBaseAddress.String(),
	}
}

// NewMsgChangeModerator returns a new MsgChangeModerator with a new moderator
func NewMsgChangeModerator(moderator sdk.AccAddress, newModerator sdk.AccAddress) *MsgChangeModerator {
	return &MsgChangeModerator{
		ModeratorAddress:    moderator.String(),
		NewModeratorAddress: newModerator.String(),
	}
}

// NewMsgResetTotalBurned returns a new MsgResetTotalBurned with a new total burned
func NewMsgResetTotalBurned(moderator sdk.AccAddress, denom string, amount math.Int) *MsgResetTotalBurned {
	return &MsgResetTotalBurned{
		ModeratorAddress: moderator.String(),
		Denom:            denom,
		Amount:           amount,
	}
}
