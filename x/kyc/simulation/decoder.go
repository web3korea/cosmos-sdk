package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/cosmos/cosmos-sdk/x/kyc/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding kyc type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.KYCStoreKeyPrefix):
			var kycA, kycB types.KYC
			cdc.MustUnmarshal(kvA.Value, &kycA)
			cdc.MustUnmarshal(kvB.Value, &kycB)
			return fmt.Sprintf("%v\n%v", kycA, kycB)
		case bytes.Equal(kvA.Key[:1], types.ValidatorStoreKeyPrefix):
			var validatorA, validatorB types.Validator
			cdc.MustUnmarshal(kvA.Value, &validatorA)
			cdc.MustUnmarshal(kvB.Value, &validatorB)
			return fmt.Sprintf("%v\n%v", validatorA, validatorB)
		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}