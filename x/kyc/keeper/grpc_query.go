package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/x/kyc/types"
)

var _ types.QueryServer = Keeper{}

// QueryGetKYC implements the Query/GetKYC gRPC method
func (k Keeper) QueryGetKYC(c context.Context, req *types.QueryGetKYCRequest) (*types.QueryGetKYCResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	kyc, found := k.GetKYCRaw(ctx, req.Address)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetKYCResponse{KYC: kyc}, nil
}

// ListKYC implements the Query/ListKYC gRPC method
func (k Keeper) ListKYC(c context.Context, req *types.QueryListKYCRequest) (*types.QueryListKYCResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	kycStore := prefix.NewStore(store, types.KYCStoreKeyPrefix)

	var kycs []types.KYC
	pageRes, err := query.Paginate(kycStore, req.Pagination, func(key []byte, value []byte) error {
		var kyc types.KYC
		if err := k.cdc.Unmarshal(value, &kyc); err != nil {
			return err
		}
		kycs = append(kycs, kyc)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListKYCResponse{KYC: kycs, Pagination: pageRes}, nil
}

// QueryGetValidator implements the Query/GetValidator gRPC method
func (k Keeper) QueryGetValidator(c context.Context, req *types.QueryGetValidatorRequest) (*types.QueryGetValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	validator, found := k.GetValidatorRaw(ctx, req.Address)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetValidatorResponse{Validator: validator}, nil
}

// ListValidators implements the Query/ListValidators gRPC method
func (k Keeper) ListValidators(c context.Context, req *types.QueryListValidatorsRequest) (*types.QueryListValidatorsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	validatorStore := prefix.NewStore(store, types.ValidatorStoreKeyPrefix)

	var validators []types.Validator
	pageRes, err := query.Paginate(validatorStore, req.Pagination, func(key []byte, value []byte) error {
		var validator types.Validator
		if err := k.cdc.Unmarshal(value, &validator); err != nil {
			return err
		}
		validators = append(validators, validator)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListValidatorsResponse{Validators: validators, Pagination: pageRes}, nil
}

// Params implements the Query/Params gRPC method
func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}