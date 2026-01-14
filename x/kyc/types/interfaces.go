package types

import (
	"context"
)

// MsgServer defines the KYC msg service.
type MsgServer interface {
	RegisterKYC(context.Context, *MsgRegisterKYC) (*MsgRegisterKYCResponse, error)
	ApproveKYC(context.Context, *MsgApproveKYC) (*MsgApproveKYCResponse, error)
	RejectKYC(context.Context, *MsgRejectKYC) (*MsgRejectKYCResponse, error)
}

// QueryServer defines the KYC query service.
type QueryServer interface {
	GetKYC(context.Context, *QueryGetKYCRequest) (*QueryGetKYCResponse, error)
	ListKYC(context.Context, *QueryListKYCRequest) (*QueryListKYCResponse, error)
	GetValidator(context.Context, *QueryGetValidatorRequest) (*QueryGetValidatorResponse, error)
	ListValidators(context.Context, *QueryListValidatorsRequest) (*QueryListValidatorsResponse, error)
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
}