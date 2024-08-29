package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/onomyprotocol/onomy/x/psm/types"
)

// QueryServer is keep wrapper which provides query capabilities.
type QueryServer struct {
	keeper Keeper
}

// NewQueryServer creates a new instance of QueryServer.
func NewQueryServer(keeper Keeper) *QueryServer {
	return &QueryServer{
		keeper: keeper,
	}
}

var _ types.QueryServer = QueryServer{}

// Params return dao module current params values.
func (q QueryServer) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: q.keeper.GetParams(ctx)}, nil
}

func (q QueryServer) Stablecoin(c context.Context, req *types.QueryStablecoinRequest) (*types.QueryStablecoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	stablecoin, found := q.keeper.GetStablecoin(ctx, req.Denom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "not found stablecoin %s", req.Denom)
	}

	moduleAddr := q.keeper.accountKeeper.GetModuleAddress(types.ModuleName)
	totalStablecoinLock := q.keeper.bankKeeper.GetBalance(ctx, moduleAddr, req.Denom).Amount

	return &types.QueryStablecoinResponse{
		Stablecoin:       stablecoin,
		CurrentTotal:     totalStablecoinLock,
		SwapableQuantity: stablecoin.LimitTotal.Sub(totalStablecoinLock),
	}, nil
}
