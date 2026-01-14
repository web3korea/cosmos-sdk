package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestQueryKycVerified(t *testing.T) {
	// Given
	suite := testdata.NewTestSuite()
	keeper := suite.Keeper.(types.AccountKeeper)
	queryServer := suite.QueryServer.(types.QueryServer)
	ctx := suite.Ctx
	account := suite.TestAccs[0]

	// Create base account with KYC verified
	baseAcc := types.NewBaseAccountWithAddress(account)
	baseAcc.SetAccountNumber(1)
	baseAcc.SetSequence(0)
	baseAcc.SetKycVerified(true)
	keeper.SetAccount(ctx, baseAcc)

	req := &types.QueryKycVerifiedRequest{
		Address: account.String(),
	}

	// When
	resp, err := queryServer.KycVerified(ctx, req)

	// Then
	require.NoError(t, err)
	require.True(t, resp.KycVerified)
}

func TestQueryKycVerified_NotVerified(t *testing.T) {
	// Given
	suite := testdata.NewTestSuite()
	keeper := suite.Keeper.(types.AccountKeeper)
	queryServer := suite.QueryServer.(types.QueryServer)
	ctx := suite.Ctx
	account := suite.TestAccs[0]

	// Create base account without KYC verified (default false)
	baseAcc := types.NewBaseAccountWithAddress(account)
	baseAcc.SetAccountNumber(1)
	baseAcc.SetSequence(0)
	keeper.SetAccount(ctx, baseAcc)

	req := &types.QueryKycVerifiedRequest{
		Address: account.String(),
	}

	// When
	resp, err := queryServer.KycVerified(ctx, req)

	// Then
	require.NoError(t, err)
	require.False(t, resp.KycVerified)
}

func TestQueryKycVerified_AccountNotFound(t *testing.T) {
	// Given
	suite := testdata.NewTestSuite()
	queryServer := suite.QueryServer.(types.QueryServer)
	ctx := suite.Ctx
	account := suite.TestAccs[0]

	req := &types.QueryKycVerifiedRequest{
		Address: account.String(),
	}

	// When
	resp, err := queryServer.KycVerified(ctx, req)

	// Then
	require.NoError(t, err)
	require.False(t, resp.KycVerified) // Default to false for non-existent accounts
}