package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterKYC{}, "kyc/RegisterKYC", nil)
	cdc.RegisterConcrete(&MsgApproveKYC{}, "kyc/ApproveKYC", nil)
	cdc.RegisterConcrete(&MsgRejectKYC{}, "kyc/RejectKYC", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterKYC{},
		&MsgApproveKYC{},
		&MsgRejectKYC{},
	)
}

var (
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)