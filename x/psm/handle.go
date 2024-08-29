package psm

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/onomyprotocol/onomy/x/psm/keeper"
	"github.com/onomyprotocol/onomy/x/psm/types"
)

// NewHandler ...
func NewHandler(k *keeper.Keeper) sdk.Handler {
	msgServer := keeper.NewMsgServerImpl(k)

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) { // nolint:gocritic //the module doesn't support messages handling
		case *types.MsgSwapToIST:
			res, err := msgServer.SwapToIST(sdk.WrapSDKContext(ctx), msg)

			fmt.Println("99999999 k	has ook")
			return sdk.WrapServiceResult(ctx, res, err)
		case *types.MsgSwapToStablecoin:
			res, err := msgServer.SwapToStablecoin(sdk.WrapSDKContext(ctx), msg)
			return sdk.WrapServiceResult(ctx, res, err)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
