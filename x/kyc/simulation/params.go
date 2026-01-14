package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/simulation"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/kyc/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(_ *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyKYCExpiryDuration),
			func(r *rand.Rand) string {
				return string(types.ModuleCdc.MustMarshalJSON(sdk.NewInt(int64(31536000 + r.Intn(31536000)))))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMinValidatorCount),
			func(r *rand.Rand) string {
				return string(types.ModuleCdc.MustMarshalJSON(sdk.NewInt(int64(1 + r.Intn(5)))))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMaxValidatorCount),
			func(r *rand.Rand) string {
				return string(types.ModuleCdc.MustMarshalJSON(sdk.NewInt(int64(5 + r.Intn(15)))))
			},
		),
	}
}