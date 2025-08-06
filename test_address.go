package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	// Set bech32 prefix to guru
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("guru", "gurupub")
	config.Seal()

	// Use the test mnemonic to generate the expected address
	// This is a simplified version - in reality you'd need to derive from the mnemonic properly
	priv := secp256k1.GenPrivKey()
	addr := sdk.AccAddress(priv.PubKey().Address()).String()
	fmt.Printf("Sample guru address: %s\n", addr)

	// Just to show the format
	fmt.Printf("Test mnemonic: %s\n", testdata.TestMnemonic)
}
