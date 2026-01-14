package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/kyc/keeper"
	"github.com/cosmos/cosmos-sdk/x/kyc/types"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	k           keeper.Keeper
	addrs       []sdk.AccAddress
	queryClient types.QueryClient
	msgServer   types.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.ctx = testdata.NewTestContext()

	// Create test addresses
	suite.addrs = make([]sdk.AccAddress, 3)
	for i := range suite.addrs {
		suite.addrs[i] = sdk.AccAddress([]byte{byte(i)})
	}

	// Create codec
	cdc := testdata.NewTestCodec()

	// Create mock account keeper
	accountKeeper := testdata.NewMockAccountKeeper(suite.addrs...)

	// Create parameter subspace
	paramsSubspace := testdata.NewMockParamSubspace()

	// Create stores
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := sdk.NewKVStoreKey(types.MemStoreKey)

	// Create context with stores
	suite.ctx = suite.ctx.WithMultiStore(sdk.NewMultiStore(
		suite.ctx.MultiStore().CacheMultiStore(),
		storeKey,
		memStoreKey,
	))

	// Create keeper
	suite.k = keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		accountKeeper,
	)

	// Create msg server
	suite.msgServer = keeper.NewMsgServerImpl(suite.k)
}

// func TestKeeperTestSuite(t *testing.T) {
// 	suite.Run(t, new(KeeperTestSuite))
// }

// func (suite *KeeperTestSuite) TestRegisterKYC() {
// 	dateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

// 	msg := &types.MsgRegisterKYC{
// 		Sender:      suite.addrs[0].String(),
// 		FullName:    "John Doe",
// 		DateOfBirth: dateOfBirth,
// 		Country:     "US",
// 		AddressInfo: "123 Main St",
// 		IDType:      "passport",
// 		IDNumber:    "P123456",
// 	}

// 	err := suite.k.RegisterKYC(suite.ctx, msg)
// 	suite.NoError(err)

// 	// Check if KYC was registered
// 	kyc, err := suite.k.GetKYC(suite.ctx, suite.addrs[0].String())
// 	suite.NoError(err)
// 	suite.Equal(suite.addrs[0].String(), kyc.Address)
// 	suite.Equal("John Doe", kyc.FullName)
// 	suite.Equal(types.StatusPending, kyc.Status)
// }

// func (suite *KeeperTestSuite) TestRegisterKYCAlreadyExists() {
// 	// First register KYC
// 	dateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

// 	msg := &types.MsgRegisterKYC{
// 		Sender:      suite.addrs[0].String(),
// 		FullName:    "John Doe",
// 		DateOfBirth: dateOfBirth,
// 		Country:     "US",
// 		AddressInfo: "123 Main St",
// 		IDType:      "passport",
// 		IDNumber:    "P123456",
// 	}

// 	err := suite.k.RegisterKYC(suite.ctx, msg)
// 	suite.NoError(err)

// 	// Try to register again - should fail
// 	err = suite.k.RegisterKYC(suite.ctx, msg)
// 	suite.Error(err)
// 	suite.Equal(types.ErrKYCAlreadyExists, err)
// }

// func (suite *KeeperTestSuite) TestApproveKYC() {
// 	// First add a validator
// 	err := suite.k.AddValidator(suite.ctx, suite.addrs[1].String(), "Validator1", []string{"approve", "reject"})
// 	suite.NoError(err)

// 	// Register KYC
// 	dateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

// 	msg := &types.MsgRegisterKYC{
// 		Sender:      suite.addrs[0].String(),
// 		FullName:    "John Doe",
// 		DateOfBirth: dateOfBirth,
// 		Country:     "US",
// 		AddressInfo: "123 Main St",
// 		IDType:      "passport",
// 		IDNumber:    "P123456",
// 	}

// 	err = suite.k.RegisterKYC(suite.ctx, msg)
// 	suite.NoError(err)

// 	// Approve KYC
// 	approveMsg := &types.MsgApproveKYC{
// 		Sender:   suite.addrs[1].String(),
// 		User:     suite.addrs[0].String(),
// 		Comments: "Approved",
// 	}

// 	err = suite.k.ApproveKYC(suite.ctx, approveMsg)
// 	suite.NoError(err)

// 	// Check if KYC was approved
// 	kyc, err := suite.k.GetKYC(suite.ctx, suite.addrs[0].String())
// 	suite.NoError(err)
// 	suite.Equal(types.StatusApproved, kyc.Status)
// 	suite.Equal(suite.addrs[1].String(), *kyc.ReviewedBy)
// 	suite.Equal("Approved", kyc.ReviewComments)
// }

// func (suite *KeeperTestSuite) TestApproveKYCUnauthorized() {
// 	// Register KYC
// 	dateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

// 	msg := &types.MsgRegisterKYC{
// 		Sender:      suite.addrs[0].String(),
// 		FullName:    "John Doe",
// 		DateOfBirth: dateOfBirth,
// 		Country:     "US",
// 		AddressInfo: "123 Main St",
// 		IDType:      "passport",
// 		IDNumber:    "P123456",
// 	}

// 	err := suite.k.RegisterKYC(suite.ctx, msg)
// 	suite.NoError(err)

// 	// Try to approve KYC with non-validator
// 	approveMsg := &types.MsgApproveKYC{
// 		Sender:   suite.addrs[2].String(), // Not a validator
// 		User:     suite.addrs[0].String(),
// 		Comments: "Approved",
// 	}

// 	err = suite.k.ApproveKYC(suite.ctx, approveMsg)
// 	suite.Error(err)
// 	suite.Equal(types.ErrUnauthorizedValidator, err)
// }

// func (suite *KeeperTestSuite) TestRejectKYC() {
// 	// First add a validator
// 	err := suite.k.AddValidator(suite.ctx, suite.addrs[1].String(), "Validator1", []string{"approve", "reject"})
// 	suite.NoError(err)

// 	// Register KYC
// 	dateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

// 	msg := &types.MsgRegisterKYC{
// 		Sender:      suite.addrs[0].String(),
// 		FullName:    "John Doe",
// 		DateOfBirth: dateOfBirth,
// 		Country:     "US",
// 		AddressInfo: "123 Main St",
// 		IDType:      "passport",
// 		IDNumber:    "P123456",
// 	}

// 	err = suite.k.RegisterKYC(suite.ctx, msg)
// 	suite.NoError(err)

// 	// Reject KYC
// 	rejectMsg := &types.MsgRejectKYC{
// 		Sender:   suite.addrs[1].String(),
// 		User:     suite.addrs[0].String(),
// 		Comments: "Invalid documents",
// 	}

// 	err = suite.k.RejectKYC(suite.ctx, rejectMsg)
// 	suite.NoError(err)

// 	// Check if KYC was rejected
// 	kyc, err := suite.k.GetKYC(suite.ctx, suite.addrs[0].String())
// 	suite.NoError(err)
// 	suite.Equal(types.StatusRejected, kyc.Status)
// 	suite.Equal(suite.addrs[1].String(), *kyc.ReviewedBy)
// 	suite.Equal("Invalid documents", kyc.ReviewComments)
// }

// func (suite *KeeperTestSuite) TestAddValidator() {
// 	err := suite.k.AddValidator(suite.ctx, suite.addrs[0].String(), "Validator1", []string{"approve", "reject"})
// 	suite.NoError(err)

// 	// Check if validator was added
// 	validator, err := suite.k.GetValidator(suite.ctx, suite.addrs[0].String())
// 	suite.NoError(err)
// 	suite.Equal(suite.addrs[0].String(), validator.Address)
// 	suite.Equal("Validator1", validator.Name)
// 	suite.Equal([]string{"approve", "reject"}, validator.Permissions)
// 	suite.True(validator.IsActive)
// }

// func (suite *KeeperTestSuite) TestAddValidatorAlreadyExists() {
// 	// Add validator first time
// 	err := suite.k.AddValidator(suite.ctx, suite.addrs[0].String(), "Validator1", []string{"approve", "reject"})
// 	suite.NoError(err)

// 	// Try to add again - should fail
// 	err = suite.k.AddValidator(suite.ctx, suite.addrs[0].String(), "Validator1", []string{"approve", "reject"})
// 	suite.Error(err)
// 	suite.Equal(types.ErrValidatorAlreadyExists, err)
// }

// func (suite *KeeperTestSuite) TestRemoveValidator() {
// 	// Add validator first
// 	err := suite.k.AddValidator(suite.ctx, suite.addrs[0].String(), "Validator1", []string{"approve", "reject"})
// 	suite.NoError(err)

// 	// Remove validator
// 	err = suite.k.RemoveValidator(suite.ctx, suite.addrs[0].String())
// 	suite.NoError(err)

// 	// Check if validator was removed
// 	validator, err := suite.k.GetValidator(suite.ctx, suite.addrs[0].String())
// 	suite.NoError(err)
// 	suite.False(validator.IsActive)
// }

// func (suite *KeeperTestSuite) TestGenesis() {
// 	// Create test data
// 	dateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
// 	kyc := types.KYC{
// 		Address:     suite.addrs[0].String(),
// 		FullName:    "John Doe",
// 		DateOfBirth: dateOfBirth,
// 		Country:     "US",
// 		AddressInfo: "123 Main St",
// 		IDType:      "passport",
// 		IDNumber:    "P123456",
// 		Status:      types.StatusApproved,
// 		SubmittedAt: time.Now(),
// 		ExpiresAt:   time.Now().AddDate(1, 0, 0),
// 	}

// 	validator := types.Validator{
// 		Address:     suite.addrs[1].String(),
// 		Name:        "Validator1",
// 		Permissions: []string{"approve", "reject"},
// 		IsActive:    true,
// 	}

// 	genesisState := types.GenesisState{
// 		Params:      types.DefaultParams(),
// 		KYCRecords:  []types.KYC{kyc},
// 		Validators:  []types.Validator{validator},
// 	}

// 	// Test InitGenesis
// 	suite.k.InitGenesis(suite.ctx, &genesisState)

// 	// Test ExportGenesis
// 	exportedGenesis := suite.k.ExportGenesis(suite.ctx)
// 	suite.Equal(genesisState.Params, exportedGenesis.Params)
// 	suite.Len(exportedGenesis.KYCRecords, 1)
// 	suite.Len(exportedGenesis.Validators, 1)
// 	suite.Equal(kyc.Address, exportedGenesis.KYCRecords[0].Address)
// 	suite.Equal(validator.Address, exportedGenesis.Validators[0].Address)
// }