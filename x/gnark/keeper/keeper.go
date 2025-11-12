package keeper

import (
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/gnark/types"
)

// Keeper manages the gnark module state.
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService store.KVStoreService
}

// NewKeeper creates a new gnark keeper instance.
func NewKeeper(cdc codec.BinaryCodec, key *storetypes.KVStoreKey) Keeper {
	if key == nil {
		panic("gnark store key is nil")
	}
	return Keeper{
		cdc:          cdc,
		storeService: runtime.NewKVStoreService(key),
	}
}

// Logger returns the module logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) kvStore(ctx sdk.Context) storetypes.KVStore {
	return runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
}

// SetCircuit persists the provided circuit.
func (k Keeper) SetCircuit(ctx sdk.Context, circuit types.Circuit) {
	store := k.kvStore(ctx)
	bz := k.cdc.MustMarshal(&circuit)
	store.Set(types.CircuitKey(circuit.Id), bz)
}

// GetCircuit fetches a circuit by id.
func (k Keeper) GetCircuit(ctx sdk.Context, id string) (types.Circuit, bool) {
	store := k.kvStore(ctx)
	bz := store.Get(types.CircuitKey(id))
	if bz == nil {
		return types.Circuit{}, false
	}

	var circuit types.Circuit
	k.cdc.MustUnmarshal(bz, &circuit)
	return circuit, true
}

// RemoveCircuit deletes the circuit from store.
func (k Keeper) RemoveCircuit(ctx sdk.Context, id string) {
	store := k.kvStore(ctx)
	store.Delete(types.CircuitKey(id))
}

// GetAllCircuits returns every stored circuit.
func (k Keeper) GetAllCircuits(ctx sdk.Context) []types.Circuit {
	store := prefix.NewStore(k.kvStore(ctx), types.CircuitKeyPrefix)
	iterator := store.Iterator(nil, nil)
	defer iterator.Close()

	var circuits []types.Circuit
	for ; iterator.Valid(); iterator.Next() {
		var circuit types.Circuit
		k.cdc.MustUnmarshal(iterator.Value(), &circuit)
		circuits = append(circuits, circuit)
	}
	return circuits
}
