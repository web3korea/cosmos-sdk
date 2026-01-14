package simulation

import (
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/cosmos/cosmos-sdk/x/kyc/types"
)

// RandomizedGenState generates a random GenesisState for the kyc module
func RandomizedGenState(simState *module.SimulationState) {
	var kycRecords []types.KYC
	var validators []types.Validator

	// Generate random KYC records
	kycCount := simtypes.RandIntBetween(simState.Rand, 0, 10)
	for i := 0; i < kycCount; i++ {
		kyc := types.KYC{
			Address:      simtypes.RandomAccounts(simState.Rand, 1)[0].Address.String(),
			FullName:     simtypes.RandStringOfLength(simState.Rand, 10),
			DateOfBirth:  time.Now().AddDate(-simtypes.RandIntBetween(simState.Rand, 18, 80), 0, 0),
			Country:      simtypes.RandStringOfLength(simState.Rand, 2),
			AddressInfo:  simtypes.RandStringOfLength(simState.Rand, 20),
			IDType:       "passport",
			IDNumber:     simtypes.RandStringOfLength(simState.Rand, 10),
			Status:       types.StatusApproved,
			SubmittedAt:  time.Now(),
			ExpiresAt:    time.Now().AddDate(1, 0, 0),
		}
		kycRecords = append(kycRecords, kyc)
	}

	// Generate random validators
	validatorCount := simtypes.RandIntBetween(simState.Rand, 1, 5)
	accounts := simtypes.RandomAccounts(simState.Rand, validatorCount)
	for _, acc := range accounts {
		validator := types.Validator{
			Address:     acc.Address.String(),
			Name:        simtypes.RandStringOfLength(simState.Rand, 8),
			Permissions: []string{"approve", "reject"},
			IsActive:    true,
		}
		validators = append(validators, validator)
	}

	genesisState := types.GenesisState{
		Params:      types.DefaultParams(),
		KYCRecords:  kycRecords,
		Validators:  validators,
	}

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&genesisState)
}