package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// KYC module errors
var (
	ErrKYCAlreadyExists      = sdkerrors.Register(ModuleName, 2, "KYC already exists for this address")
	ErrKYCNotFound           = sdkerrors.Register(ModuleName, 3, "KYC not found")
	ErrKYCNotApproved        = sdkerrors.Register(ModuleName, 4, "KYC not approved")
	ErrKYCExpired            = sdkerrors.Register(ModuleName, 5, "KYC expired")
	ErrKYCAlreadyApproved    = sdkerrors.Register(ModuleName, 6, "KYC already approved")
	ErrKYCAlreadyRejected    = sdkerrors.Register(ModuleName, 7, "KYC already rejected")
	ErrUnauthorizedValidator = sdkerrors.Register(ModuleName, 8, "unauthorized validator")
	ErrValidatorNotFound     = sdkerrors.Register(ModuleName, 9, "validator not found")
	ErrValidatorAlreadyExists = sdkerrors.Register(ModuleName, 10, "validator already exists")
	ErrInvalidValidatorAddress = sdkerrors.Register(ModuleName, 11, "invalid validator address")
	ErrInsufficientValidators  = sdkerrors.Register(ModuleName, 12, "insufficient validators")
	ErrInvalidKYCStatus      = sdkerrors.Register(ModuleName, 13, "invalid KYC status")
	ErrInvalidPermissions    = sdkerrors.Register(ModuleName, 14, "invalid permissions")
)