package types_test

import (
	"github.com/cosmos/cosmos-sdk/simulateapp"
)

var (
	ecdc                  = simulateapp.MakeTestEncodingConfig()
	appCodec, legacyAmino = ecdc.Codec, ecdc.Amino
)
