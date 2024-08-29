package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/psm/types"
)

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.ps.GetParamSet(ctx, &params)
	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.ps.SetParamSet(ctx, &params)
}

func (k Keeper) AcceptablePriceRatio(ctx sdk.Context) (res int64) {
	k.ps.Get(ctx, types.KeyAcceptablePriceRatio, &res)
	return
}

func (k Keeper) TotalLimit(ctx sdk.Context) (res sdk.Int) {
	k.ps.Get(ctx, types.KeyLimmitTotal, &res)
	return
}
