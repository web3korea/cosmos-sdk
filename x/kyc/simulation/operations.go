package simulation

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/cosmos/cosmos-sdk/x/kyc/keeper"
	"github.com/cosmos/cosmos-sdk/x/kyc/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgRegisterKYC = "op_weight_msg_register_kyc"
	OpWeightMsgApproveKYC  = "op_weight_msg_approve_kyc"
	OpWeightMsgRejectKYC   = "op_weight_msg_reject_kyc"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONCodec,
	ak types.AccountKeeper, k keeper.Keeper,
) []simtypes.WeightedOperation {
	return []simtypes.WeightedOperation{
		simulation.NewWeightedOperation(
			appParams.GetOrGenerate(cdc, OpWeightMsgRegisterKYC, nil, func(_ *rand.Rand) {
				simValue := 80
				return simValue
			}),
			SimulateMsgRegisterKYC(ak, k),
		),
		simulation.NewWeightedOperation(
			appParams.GetOrGenerate(cdc, OpWeightMsgApproveKYC, nil, func(_ *rand.Rand) {
				simValue := 50
				return simValue
			}),
			SimulateMsgApproveKYC(ak, k),
		),
		simulation.NewWeightedOperation(
			appParams.GetOrGenerate(cdc, OpWeightMsgRejectKYC, nil, func(_ *rand.Rand) {
				simValue := 30
				return simValue
			}),
			SimulateMsgRejectKYC(ak, k),
		),
	}
}

// SimulateMsgRegisterKYC simulates a MsgRegisterKYC message
func SimulateMsgRegisterKYC(ak types.AccountKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		// Check if KYC already exists
		if k.HasKYC(ctx, simAccount.Address.String()) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRegisterKYC, "KYC already exists"), nil, nil
		}

		msg := types.MsgRegisterKYC{
			Sender:      simAccount.Address.String(),
			FullName:    simtypes.RandStringOfLength(r, 10),
			DateOfBirth: time.Now().AddDate(-simtypes.RandIntBetween(r, 18, 80), 0, 0),
			Country:     simtypes.RandStringOfLength(r, 2),
			AddressInfo: simtypes.RandStringOfLength(r, 20),
			IDType:      "passport",
			IDNumber:    simtypes.RandStringOfLength(r, 10),
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simtypes.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      nil,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgApproveKYC simulates a MsgApproveKYC message
func SimulateMsgApproveKYC(ak types.AccountKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		// Check if sender is a validator
		if !k.IsValidator(ctx, simAccount.Address.String()) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgApproveKYC, "not a validator"), nil, nil
		}

		// Get a random KYC record that can be approved
		allKYC := k.GetAllKYC(ctx)
		if len(allKYC) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgApproveKYC, "no KYC records"), nil, nil
		}

		kyc := allKYC[r.Intn(len(allKYC))]
		if kyc.Status != types.StatusPending {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgApproveKYC, "KYC not pending"), nil, nil
		}

		msg := types.MsgApproveKYC{
			Sender:   simAccount.Address.String(),
			User:     kyc.Address,
			Comments: "Approved via simulation",
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simtypes.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      nil,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgRejectKYC simulates a MsgRejectKYC message
func SimulateMsgRejectKYC(ak types.AccountKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)

		// Check if sender is a validator
		if !k.IsValidator(ctx, simAccount.Address.String()) {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRejectKYC, "not a validator"), nil, nil
		}

		// Get a random KYC record that can be rejected
		allKYC := k.GetAllKYC(ctx)
		if len(allKYC) == 0 {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRejectKYC, "no KYC records"), nil, nil
		}

		kyc := allKYC[r.Intn(len(allKYC))]
		if kyc.Status != types.StatusPending {
			return simtypes.NoOpMsg(types.ModuleName, types.TypeMsgRejectKYC, "KYC not pending"), nil, nil
		}

		msg := types.MsgRejectKYC{
			Sender:   simAccount.Address.String(),
			User:     kyc.Address,
			Comments: "Rejected via simulation",
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           simtypes.MakeTestEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      nil,
			ModuleName:      types.ModuleName,
		}

		return simulation.GenAndDeliverTxWithRandFees(txCtx)
	}
}