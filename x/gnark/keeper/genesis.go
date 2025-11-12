package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/gnark/types"
)

// InitGenesis initializes module state from genesis data.
func (k Keeper) InitGenesis(ctx sdk.Context, state *types.GenesisState) {
	if state == nil {
		state = types.DefaultGenesis()
	}

	for _, circuit := range state.Circuits {
		k.SetCircuit(ctx, circuit)
	}
}

// ExportGenesis exports the current module state.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:   types.DefaultGenesisParams(),
		Circuits: k.GetAllCircuits(ctx),
	}
}
