package types_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/x/kyc/types"
)

func TestNewKYC(t *testing.T) {
	dateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	kyc := types.NewKYC(
		"cosmos1test",
		"John Doe",
		dateOfBirth,
		"US",
		"123 Main St",
		"passport",
		"P123456",
	)

	require.Equal(t, "cosmos1test", kyc.Address)
	require.Equal(t, "John Doe", kyc.FullName)
	require.Equal(t, dateOfBirth, kyc.DateOfBirth)
	require.Equal(t, "US", kyc.Country)
	require.Equal(t, "123 Main St", kyc.AddressInfo)
	require.Equal(t, "passport", kyc.IDType)
	require.Equal(t, "P123456", kyc.IDNumber)
	require.Equal(t, types.StatusPending, kyc.Status)
	require.True(t, kyc.SubmittedAt.After(time.Now().Add(-time.Second)))
	require.True(t, kyc.ExpiresAt.After(time.Now().AddDate(0, 11, 25))) // Almost 1 year from now
}

func TestKYCIsExpired(t *testing.T) {
	// Test not expired
	kyc := types.KYC{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	require.False(t, kyc.IsExpired())

	// Test expired
	kyc = types.KYC{
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	require.True(t, kyc.IsExpired())
}

func TestKYCIsValid(t *testing.T) {
	// Test valid KYC (approved and not expired)
	kyc := types.KYC{
		Status:    types.StatusApproved,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	require.True(t, kyc.IsValid())

	// Test invalid - rejected
	kyc = types.KYC{
		Status:    types.StatusRejected,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	require.False(t, kyc.IsValid())

	// Test invalid - expired
	kyc = types.KYC{
		Status:    types.StatusApproved,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	require.False(t, kyc.IsValid())

	// Test invalid - pending
	kyc = types.KYC{
		Status:    types.StatusPending,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	require.False(t, kyc.IsValid())
}

func TestKYCValidate(t *testing.T) {
	dateOfBirth := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	// Test valid KYC
	kyc := types.KYC{
		Address:     "cosmos1test",
		FullName:    "John Doe",
		DateOfBirth: dateOfBirth,
		Country:     "US",
		IDType:      "passport",
		IDNumber:    "P123456",
	}
	err := kyc.Validate()
	require.NoError(t, err)

	// Test invalid - empty address
	kyc.Address = ""
	err = kyc.Validate()
	require.Error(t, err)

	// Test invalid - empty full name
	kyc.Address = "cosmos1test"
	kyc.FullName = ""
	err = kyc.Validate()
	require.Error(t, err)

	// Test invalid - empty country
	kyc.FullName = "John Doe"
	kyc.Country = ""
	err = kyc.Validate()
	require.Error(t, err)

	// Test invalid - empty ID type
	kyc.Country = "US"
	kyc.IDType = ""
	err = kyc.Validate()
	require.Error(t, err)

	// Test invalid - empty ID number
	kyc.IDType = "passport"
	kyc.IDNumber = ""
	err = kyc.Validate()
	require.Error(t, err)
}

func TestKYCUpdateStatus(t *testing.T) {
	kyc := types.KYC{
		Status: types.StatusPending,
	}

	// Update to approved
	kyc.UpdateStatus(types.StatusApproved, "validator1", "Approved")
	require.Equal(t, types.StatusApproved, kyc.Status)
	require.Equal(t, "validator1", *kyc.ReviewedBy)
	require.Equal(t, "Approved", kyc.ReviewComments)
	require.NotNil(t, kyc.ReviewedAt)
	require.True(t, kyc.ReviewedAt.After(time.Now().Add(-time.Second)))
}

func TestKYCAddDocument(t *testing.T) {
	kyc := types.KYC{
		Documents: []types.KYCDocument{},
	}

	kyc.AddDocument("passport", "ipfs://hash", "doc_hash_123")

	require.Len(t, kyc.Documents, 1)
	doc := kyc.Documents[0]
	require.Equal(t, "passport", doc.Type)
	require.Equal(t, "ipfs://hash", doc.URL)
	require.Equal(t, "doc_hash_123", doc.Hash)
	require.True(t, doc.UploadedAt.After(time.Now().Add(-time.Second)))
}