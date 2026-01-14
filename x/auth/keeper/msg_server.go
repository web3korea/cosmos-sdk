package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the auth MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) SetKycVerified(goCtx context.Context, msg *types.MsgSetKycVerified) (*types.MsgSetKycVerifiedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	admin, err := sdk.AccAddressFromBech32(msg.Admin)
	if err != nil {
		return nil, err
	}

	account, err := sdk.AccAddressFromBech32(msg.Account)
	if err != nil {
		return nil, err
	}

	// Check if the admin has permission to set KYC status
	if !k.IsAdmin(ctx, admin) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "only admin can set KYC verification status")
	}

	// Get the account
	acc := k.GetAccount(ctx, account)
	if acc == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNotFound, "account not found")
	}

	// Update KYC verification status
	baseAcc, ok := acc.(*types.BaseAccount)
	if !ok {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidType, "account is not a BaseAccount")
	}

	err = baseAcc.SetKycVerified(msg.KycVerified)
	if err != nil {
		return nil, err
	}

	// Save the updated account
	k.SetAccount(ctx, baseAcc)

	return &types.MsgSetKycVerifiedResponse{}, nil
}