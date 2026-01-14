package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/kyc/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the KYC MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// RegisterKYC implements types.MsgServer
func (k msgServer) RegisterKYC(goCtx context.Context, msg *types.MsgRegisterKYC) (*types.MsgRegisterKYCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.Keeper.RegisterKYC(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgRegisterKYCResponse{}, nil
}

// ApproveKYC implements types.MsgServer
func (k msgServer) ApproveKYC(goCtx context.Context, msg *types.MsgApproveKYC) (*types.MsgApproveKYCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.Keeper.ApproveKYC(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgApproveKYCResponse{}, nil
}

// RejectKYC implements types.MsgServer
func (k msgServer) RejectKYC(goCtx context.Context, msg *types.MsgRejectKYC) (*types.MsgRejectKYCResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.Keeper.RejectKYC(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgRejectKYCResponse{}, nil
}