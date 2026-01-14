package types

const (
	// ModuleName defines the module name
	ModuleName = "kyc"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_kyc"
)

// KYC status constants
const (
	StatusPending   = "pending"
	StatusApproved  = "approved"
	StatusRejected  = "rejected"
	StatusExpired   = "expired"
	StatusSuspended = "suspended"
)

// KYC store keys
var (
	// KYCStoreKeyPrefix is the prefix for KYC store keys
	KYCStoreKeyPrefix = []byte{0x01}

	// ValidatorStoreKeyPrefix is the prefix for validator store keys
	ValidatorStoreKeyPrefix = []byte{0x02}

	// ParamsKey is the key for parameters
	ParamsKey = []byte{0x03}
)

// GetKYCKey returns the key for a specific KYC record
func GetKYCKey(address string) []byte {
	return append(KYCStoreKeyPrefix, []byte(address)...)
}

// GetValidatorKey returns the key for a validator
func GetValidatorKey(address string) []byte {
	return append(ValidatorStoreKeyPrefix, []byte(address)...)
}