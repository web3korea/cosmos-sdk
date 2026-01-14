package types

import (
	"github.com/cosmos/cosmos-sdk/types/query"
)

// Query endpoints supported by the kyc querier
const (
	QueryGetKYC     = "get_kyc"
	QueryListKYC    = "list_kyc"
	QueryGetValidator = "get_validator"
	QueryListValidators = "list_validators"
	QueryParams     = "params"
)

// QueryGetKYCRequest is the request type for the Query/GetKYC RPC method
type QueryGetKYCRequest struct {
	Address string `json:"address,omitempty" yaml:"address"`
}

// QueryGetKYCResponse is the response type for the Query/GetKYC RPC method
type QueryGetKYCResponse struct {
	KYC KYC `json:"kyc,omitempty" yaml:"kyc"`
}

// QueryListKYCRequest is the request type for the Query/ListKYC RPC method
type QueryListKYCRequest struct {
	Pagination *query.PageRequest `json:"pagination,omitempty" yaml:"pagination"`
}

// QueryListKYCResponse is the response type for the Query/ListKYC RPC method
type QueryListKYCResponse struct {
	KYC        []KYC             `json:"kyc,omitempty" yaml:"kyc"`
	Pagination *query.PageResponse `json:"pagination,omitempty" yaml:"pagination"`
}

// QueryGetValidatorRequest is the request type for the Query/GetValidator RPC method
type QueryGetValidatorRequest struct {
	Address string `json:"address,omitempty" yaml:"address"`
}

// QueryGetValidatorResponse is the response type for the Query/GetValidator RPC method
type QueryGetValidatorResponse struct {
	Validator Validator `json:"validator,omitempty" yaml:"validator"`
}

// QueryListValidatorsRequest is the request type for the Query/ListValidators RPC method
type QueryListValidatorsRequest struct {
	Pagination *query.PageRequest `json:"pagination,omitempty" yaml:"pagination"`
}

// QueryListValidatorsResponse is the response type for the Query/ListValidators RPC method
type QueryListValidatorsResponse struct {
	Validators []Validator         `json:"validators,omitempty" yaml:"validators"`
	Pagination *query.PageResponse `json:"pagination,omitempty" yaml:"pagination"`
}

// QueryParamsRequest is the request type for the Query/Params RPC method
type QueryParamsRequest struct{}

// QueryParamsResponse is the response type for the Query/Params RPC method
type QueryParamsResponse struct {
	Params Params `json:"params,omitempty" yaml:"params"`
}