package types

import (
	"time"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// KYC represents a KYC record
type KYC struct {
	// Address of the user
	Address string `json:"address" yaml:"address"`

	// Personal information
	FullName    string    `json:"full_name" yaml:"full_name"`
	DateOfBirth time.Time `json:"date_of_birth" yaml:"date_of_birth"`
	Country     string    `json:"country" yaml:"country"`
	AddressInfo string    `json:"address_info" yaml:"address_info"`

	// Identification documents
	IDType   string `json:"id_type" yaml:"id_type"`   // passport, drivers_license, national_id
	IDNumber string `json:"id_number" yaml:"id_number"`

	// KYC status and timestamps
	Status           string    `json:"status" yaml:"status"`
	SubmittedAt      time.Time `json:"submitted_at" yaml:"submitted_at"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty" yaml:"reviewed_at,omitempty"`
	ReviewedBy       string     `json:"reviewed_by,omitempty" yaml:"reviewed_by,omitempty"`
	ReviewComments   string     `json:"review_comments,omitempty" yaml:"review_comments,omitempty"`
	ExpiresAt        time.Time `json:"expires_at" yaml:"expires_at"`

	// Additional metadata
	Documents []KYCDocument `json:"documents" yaml:"documents"`
}

// KYCDocument represents a document associated with KYC
type KYCDocument struct {
	Type        string `json:"type" yaml:"type"`               // photo, passport_scan, etc.
	URL         string `json:"url" yaml:"url"`                // IPFS hash or URL
	Hash        string `json:"hash" yaml:"hash"`              // Document hash for verification
	UploadedAt  time.Time `json:"uploaded_at" yaml:"uploaded_at"`
}

// Validator represents a KYC validator
type Validator struct {
	Address     string   `json:"address" yaml:"address"`
	Name        string   `json:"name" yaml:"name"`
	Permissions []string `json:"permissions" yaml:"permissions"` // approve, reject, suspend
	IsActive    bool     `json:"is_active" yaml:"is_active"`
}

// NewKYC creates a new KYC record
func NewKYC(address string, fullName string, dateOfBirth time.Time, country, addressInfo, idType, idNumber string) KYC {
	now := time.Now()
	expiresAt := now.AddDate(1, 0, 0) // Expires in 1 year

	return KYC{
		Address:      address,
		FullName:     fullName,
		DateOfBirth:  dateOfBirth,
		Country:      country,
		AddressInfo:  addressInfo,
		IDType:       idType,
		IDNumber:     idNumber,
		Status:       StatusPending,
		SubmittedAt:  now,
		ExpiresAt:    expiresAt,
		Documents:    []KYCDocument{},
	}
}

// IsExpired checks if the KYC record is expired
func (k KYC) IsExpired() bool {
	return time.Now().After(k.ExpiresAt)
}

// IsValid checks if the KYC record is valid (approved and not expired)
func (k KYC) IsValid() bool {
	return k.Status == StatusApproved && !k.IsExpired()
}

// Validate validates the KYC record
func (k KYC) Validate() error {
	if k.Address == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "address cannot be empty")
	}
	if k.FullName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "full name cannot be empty")
	}
	if k.Country == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "country cannot be empty")
	}
	if k.IDType == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "ID type cannot be empty")
	}
	if k.IDNumber == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "ID number cannot be empty")
	}
	return nil
}

// UpdateStatus updates the KYC status
func (k *KYC) UpdateStatus(status string, reviewedBy string, comments string) {
	now := time.Now()
	k.Status = status
	k.ReviewedAt = &now
	k.ReviewedBy = reviewedBy
	k.ReviewComments = comments
}

// AddDocument adds a document to the KYC record
func (k *KYC) AddDocument(docType, url, hash string) {
	doc := KYCDocument{
		Type:       docType,
		URL:        url,
		Hash:       hash,
		UploadedAt: time.Now(),
	}
	k.Documents = append(k.Documents, doc)
}