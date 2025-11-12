package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/onomyprotocol/onomy/x/gnark/types"
)

var _ types.QueryServer = Keeper{}

// Circuit returns a single stored circuit by id.
func (k Keeper) Circuit(ctx context.Context, req *types.QueryCircuitRequest) (*types.QueryCircuitResponse, error) {
	if req == nil || strings.TrimSpace(req.CircuitId) == "" {
		return nil, sdkerrors.ErrInvalidRequest
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	circuit, found := k.GetCircuit(sdkCtx, req.CircuitId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrCircuitNotFound, "circuit %s", req.CircuitId)
	}

	circuitCopy := circuit
	return &types.QueryCircuitResponse{Circuit: &circuitCopy}, nil
}

// Circuits lists all stored circuits with pagination.
func (k Keeper) Circuits(ctx context.Context, req *types.QueryCircuitsRequest) (*types.QueryCircuitsResponse, error) {
	if req == nil {
		return nil, sdkerrors.ErrInvalidRequest
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	store := prefix.NewStore(k.kvStore(sdkCtx), types.CircuitKeyPrefix)

	var circuits []types.Circuit
	pageRes, err := query.Paginate(store, req.Pagination, func(_ []byte, value []byte) error {
		var circuit types.Circuit
		k.cdc.MustUnmarshal(value, &circuit)
		circuits = append(circuits, circuit)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &types.QueryCircuitsResponse{
		Circuits:   circuits,
		Pagination: pageRes,
	}, nil
}
