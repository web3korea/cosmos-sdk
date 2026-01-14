package types_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestMsgSetKycVerified(t *testing.T) {
	// Given
	suite := testdata.NewTestSuite()
	keeper := suite.Keeper.(types.AccountKeeper)
	ctx := suite.Ctx
	admin := suite.TestAccs[0]
	targetAccount := suite.TestAccs[1]

	// Create base account for target
	baseAcc := types.NewBaseAccountWithAddress(targetAccount)
	baseAcc.SetAccountNumber(1)
	baseAcc.SetSequence(0)
	keeper.SetAccount(ctx, baseAcc)

	msg := &types.MsgSetKycVerified{
		Admin:        admin.String(),
		Account:      targetAccount.String(),
		KycVerified: true,
	}

	msgServer := types.NewMsgServerImpl(keeper)

	// When
	_, err := msgServer.SetKycVerified(ctx, msg)

	// Then
	require.NoError(t, err)

	// Verify the account was updated
	updatedAcc := keeper.GetAccount(ctx, targetAccount)
	require.NotNil(t, updatedAcc)
	baseUpdatedAcc, ok := updatedAcc.(*types.BaseAccount)
	require.True(t, ok)
	require.True(t, baseUpdatedAcc.GetKycVerified())
}

func TestMsgSetKycVerified_Unauthorized(t *testing.T) {
	// Given
	suite := testdata.NewTestSuite()
	keeper := suite.Keeper.(types.AccountKeeper)
	ctx := suite.Ctx
	nonAdmin := suite.TestAccs[0]
	targetAccount := suite.TestAccs[1]

	// Create base account for target
	baseAcc := types.NewBaseAccountWithAddress(targetAccount)
	keeper.SetAccount(ctx, baseAcc)

	msg := &types.MsgSetKycVerified{
		Admin:        nonAdmin.String(), // Not an admin
		Account:      targetAccount.String(),
		KycVerified: true,
	}

	msgServer := types.NewMsgServerImpl(keeper)

	// When
	_, err := msgServer.SetKycVerified(ctx, msg)

	// Then
	require.Error(t, err)
	require.Contains(t, err.Error(), "unauthorized")
}