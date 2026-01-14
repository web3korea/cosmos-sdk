package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMsgSetKycVerified creates a new MsgSetKycVerified instance
func NewMsgSetKycVerified(admin, account sdk.AccAddress, kycVerified bool) *MsgSetKycVerified {
	return &MsgSetKycVerified{
		Admin:        admin.String(),
		Account:      account.String(),
		KycVerified: kycVerified,
	}
}