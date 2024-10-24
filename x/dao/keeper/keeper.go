// Package keeper contains dao module keeper.
package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/dao/types"
)

type (
	// Keeper is a dao keeper struct.
	Keeper struct {
		cdc codec.BinaryCodec

		storeService store.KVStoreService

		authority string

		ps types.ParamSubspace

		bankKeeper         types.BankKeeper
		accountKeeper      types.AccountKeeper
		distributionKeeper types.DistributionKeeper
		govKeeper          types.GovKeeper
		mintKeeper         types.MintKeeper
		stakingKeeper      types.StakingKeeper
	}
)

// NewKeeper creates new dao keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	authority string,
	ps types.ParamSubspace,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	distributionKeeper types.DistributionKeeper,
	govKeeper types.GovKeeper,
	mintKeeper types.MintKeeper,
	stakingKeeper types.StakingKeeper,
) *Keeper {
	// set KeyTable if it has not already been set.
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	// ensure dao module account is set.
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return &Keeper{
		cdc:                cdc,
		storeService:       storeService,
		authority:          authority,
		ps:                 ps,
		bankKeeper:         bankKeeper,
		accountKeeper:      accountKeeper,
		distributionKeeper: distributionKeeper,
		govKeeper:          govKeeper,
		mintKeeper:         mintKeeper,
		stakingKeeper:      stakingKeeper,
	}
}

// Logger returns keeper logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
