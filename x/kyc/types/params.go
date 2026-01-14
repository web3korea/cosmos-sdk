package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys
var (
	KeyKYCExpiryDuration = []byte("KYCExpiryDuration")
	KeyMinValidatorCount = []byte("MinValidatorCount")
	KeyMaxValidatorCount = []byte("MaxValidatorCount")
)

// ParamTable for kyc module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Params defines the parameters for the kyc module
type Params struct {
	// KYC expiry duration in seconds (default: 1 year)
	KYCExpiryDuration int64 `json:"kyc_expiry_duration" yaml:"kyc_expiry_duration"`

	// Minimum number of validators required
	MinValidatorCount uint32 `json:"min_validator_count" yaml:"min_validator_count"`

	// Maximum number of validators allowed
	MaxValidatorCount uint32 `json:"max_validator_count" yaml:"max_validator_count"`
}

// NewParams creates a new Params instance
func NewParams(kycExpiryDuration int64, minValidatorCount, maxValidatorCount uint32) Params {
	return Params{
		KYCExpiryDuration: kycExpiryDuration,
		MinValidatorCount: minValidatorCount,
		MaxValidatorCount: maxValidatorCount,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		31536000, // 1 year in seconds
		1,        // minimum 1 validator
		10,       // maximum 10 validators
	)
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyKYCExpiryDuration, &p.KYCExpiryDuration, validateKYCExpiryDuration),
		paramtypes.NewParamSetPair(KeyMinValidatorCount, &p.MinValidatorCount, validateMinValidatorCount),
		paramtypes.NewParamSetPair(KeyMaxValidatorCount, &p.MaxValidatorCount, validateMaxValidatorCount),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateKYCExpiryDuration(p.KYCExpiryDuration); err != nil {
		return err
	}
	if err := validateMinValidatorCount(p.MinValidatorCount); err != nil {
		return err
	}
	if err := validateMaxValidatorCount(p.MaxValidatorCount); err != nil {
		return err
	}
	if p.MinValidatorCount > p.MaxValidatorCount {
		return fmt.Errorf("min_validator_count cannot be greater than max_validator_count")
	}
	return nil
}

// validateKYCExpiryDuration validates the KYC expiry duration parameter
func validateKYCExpiryDuration(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v <= 0 {
		return fmt.Errorf("kyc_expiry_duration must be positive: %d", v)
	}
	return nil
}

// validateMinValidatorCount validates the minimum validator count parameter
func validateMinValidatorCount(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("min_validator_count must be at least 1: %d", v)
	}
	return nil
}

// validateMaxValidatorCount validates the maximum validator count parameter
func validateMaxValidatorCount(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("max_validator_count must be at least 1: %d", v)
	}
	return nil
}