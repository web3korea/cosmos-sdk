package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
)

// GenesisState defines the kyc module's genesis state
type GenesisState struct {
	// Parameters for the module
	Params Params `json:"params" yaml:"params"`

	// KYC records
	KYCRecords []KYC `json:"kyc_records" yaml:"kyc_records"`

	// Validators
	Validators []Validator `json:"validators" yaml:"validators"`
}

// NewGenesisState creates a new genesis state
func NewGenesisState(params Params, kycRecords []KYC, validators []Validator) *GenesisState {
	return &GenesisState{
		Params:     params,
		KYCRecords: kycRecords,
		Validators: validators,
	}
}

// DefaultGenesisState returns the default genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		KYCRecords: []KYC{},
		Validators: []Validator{},
	}
}

// ValidateGenesis validates the genesis state
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	// Validate KYC records
	for _, kyc := range data.KYCRecords {
		if err := kyc.Validate(); err != nil {
			return err
		}
	}

	// Validate validators
	for _, validator := range data.Validators {
		if validator.Address == "" {
			return ErrInvalidValidatorAddress
		}
	}

	return nil
}

// GetGenesisStateFromAppState returns GenesisState given raw application genesis state
func GetGenesisStateFromAppState(cdc *codec.LegacyAmino, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}