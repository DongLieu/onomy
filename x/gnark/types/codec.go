package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the legacy Amino messages.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateCircuit{}, "gnark/CreateCircuit", nil)
	cdc.RegisterConcrete(&MsgVerifyProof{}, "gnark/VerifyProof", nil)
}

// RegisterInterfaces registers the module interface types.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateCircuit{},
		&MsgVerifyProof{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	// ModuleCdc references the global module codec.
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry()) //nolint:gochecknoglobals // cosmos sdk module style
)
