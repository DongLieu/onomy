package psm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/psm/keeper"
)

// BeginBlocker does any custom logic for the DAO upon `BeginBlocker`
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	// update price and feeIn feeOut

}
