package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// KYC messages
const (
	TypeMsgRegisterKYC    = "register_kyc"
	TypeMsgUpdateKYC      = "update_kyc"
	TypeMsgApproveKYC     = "approve_kyc"
	TypeMsgRejectKYC      = "reject_kyc"
	TypeMsgSuspendKYC     = "suspend_kyc"
	TypeMsgAddValidator   = "add_validator"
	TypeMsgRemoveValidator = "remove_validator"
)

// MsgRegisterKYC represents a message to register KYC
type MsgRegisterKYC struct {
	Sender       string    `json:"sender" yaml:"sender"`
	FullName     string    `json:"full_name" yaml:"full_name"`
	DateOfBirth  time.Time `json:"date_of_birth" yaml:"date_of_birth"`
	Country      string    `json:"country" yaml:"country"`
	AddressInfo  string    `json:"address_info" yaml:"address_info"`
	IDType       string    `json:"id_type" yaml:"id_type"`
	IDNumber     string    `json:"id_number" yaml:"id_number"`
}

// ProtoMessage implements the proto.Message interface
func (msg MsgRegisterKYC) ProtoMessage() {}

// Reset implements the proto.Message interface
func (msg *MsgRegisterKYC) Reset() {}

// String implements the proto.Message interface
func (msg MsgRegisterKYC) String() string {
	return msg.Sender
}

// NewMsgRegisterKYC creates a new MsgRegisterKYC
func NewMsgRegisterKYC(sender, fullName string, dateOfBirth time.Time, country, addressInfo, idType, idNumber string) *MsgRegisterKYC {
	return &MsgRegisterKYC{
		Sender:      sender,
		FullName:    fullName,
		DateOfBirth: dateOfBirth,
		Country:     country,
		AddressInfo: addressInfo,
		IDType:      idType,
		IDNumber:    idNumber,
	}
}

// Route returns the route for this message
func (msg MsgRegisterKYC) Route() string {
	return RouterKey
}

// Type returns the type of this message
func (msg MsgRegisterKYC) Type() string {
	return TypeMsgRegisterKYC
}

// GetSigners returns the signers of this message
func (msg MsgRegisterKYC) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes returns the sign bytes for this message
func (msg MsgRegisterKYC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation
func (msg MsgRegisterKYC) ValidateBasic() error {
	if msg.Sender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if msg.FullName == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "full name cannot be empty")
	}
	if msg.Country == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "country cannot be empty")
	}
	if msg.IDType == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "ID type cannot be empty")
	}
	if msg.IDNumber == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "ID number cannot be empty")
	}
	return nil
}

// MsgApproveKYC represents a message to approve KYC
type MsgApproveKYC struct {
	Sender    string `json:"sender" yaml:"sender"`
	User      string `json:"user" yaml:"user"`
	Comments  string `json:"comments,omitempty" yaml:"comments,omitempty"`
}

// NewMsgApproveKYC creates a new MsgApproveKYC
func NewMsgApproveKYC(sender, user, comments string) *MsgApproveKYC {
	return &MsgApproveKYC{
		Sender:   sender,
		User:     user,
		Comments: comments,
	}
}

// Route returns the route for this message
func (msg MsgApproveKYC) Route() string {
	return RouterKey
}

// Type returns the type of this message
func (msg MsgApproveKYC) Type() string {
	return TypeMsgApproveKYC
}

// GetSigners returns the signers of this message
func (msg MsgApproveKYC) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes returns the sign bytes for this message
func (msg MsgApproveKYC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation
func (msg MsgApproveKYC) ValidateBasic() error {
	if msg.Sender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if msg.User == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "user cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.User); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	return nil
}

// MsgRejectKYC represents a message to reject KYC
type MsgRejectKYC struct {
	Sender    string `json:"sender" yaml:"sender"`
	User      string `json:"user" yaml:"user"`
	Comments  string `json:"comments" yaml:"comments"`
}

// NewMsgRejectKYC creates a new MsgRejectKYC
func NewMsgRejectKYC(sender, user, comments string) *MsgRejectKYC {
	return &MsgRejectKYC{
		Sender:   sender,
		User:     user,
		Comments: comments,
	}
}

// Route returns the route for this message
func (msg MsgRejectKYC) Route() string {
	return RouterKey
}

// Type returns the type of this message
func (msg MsgRejectKYC) Type() string {
	return TypeMsgRejectKYC
}

// GetSigners returns the signers of this message
func (msg MsgRejectKYC) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// GetSignBytes returns the sign bytes for this message
func (msg MsgRejectKYC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic validation
func (msg MsgRejectKYC) ValidateBasic() error {
	if msg.Sender == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "sender cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if msg.User == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "user cannot be empty")
	}
	if _, err := sdk.AccAddressFromBech32(msg.User); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, err.Error())
	}
	if msg.Comments == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "comments cannot be empty")
	}
	return nil
}

// Message response types

// MsgRegisterKYCResponse defines the response type for MsgRegisterKYC
type MsgRegisterKYCResponse struct{}

// MsgApproveKYCResponse defines the response type for MsgApproveKYC
type MsgApproveKYCResponse struct{}

// MsgRejectKYCResponse defines the response type for MsgRejectKYC
type MsgRejectKYCResponse struct{}