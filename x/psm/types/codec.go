package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&AddStableCoinProposal{}, fmt.Sprintf("%s/%s", ModuleName, ProposalTypeAddStableCoinProposal), nil)
	cdc.RegisterConcrete(&UpdatesStableCoinProposal{}, fmt.Sprintf("%s/%s", ModuleName, ProposalTypeUpdatesStableCoinProposal), nil)
}

// RegisterInterfaces registers the cdctypes interface.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&AddStableCoinProposal{},
		&UpdatesStableCoinProposal{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSwapToIST{},
		&MsgSwapToStablecoin{},
	)
}
