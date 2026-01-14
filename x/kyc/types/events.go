package types

// KYC module event types
const (
	EventTypeKYCRegistered   = "kyc_registered"
	EventTypeKYCApproved     = "kyc_approved"
	EventTypeKYCRejected     = "kyc_rejected"
	EventTypeKYCUpdated      = "kyc_updated"
	EventTypeValidatorAdded  = "validator_added"
	EventTypeValidatorRemoved = "validator_removed"
)

// KYC module event attributes
const (
	AttributeKeyAddress   = "address"
	AttributeKeyStatus    = "status"
	AttributeKeyValidator = "validator"
	AttributeKeyName      = "name"
	AttributeKeyUser      = "user"
)