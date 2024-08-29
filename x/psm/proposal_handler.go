package psm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/onomyprotocol/onomy/x/psm/cli"
	"github.com/onomyprotocol/onomy/x/psm/keeper"
	"github.com/onomyprotocol/onomy/x/psm/types"
)

var (
	AddStableCoinProposalHandler     = govclient.NewProposalHandler(cli.NewCmdSubmitAddStableCoinProposal, cli.ProposalRESTAddHandler)
	UpdatesStableCoinProposalHandler = govclient.NewProposalHandler(cli.NewCmdUpdatesStableCoinProposal, cli.ProposalRESTUpdateHandler)
)

func NewStablecoinProposalHandler(k *keeper.Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.AddStableCoinProposal:
			return k.AddStableCoinProposal(ctx, c)
		case *types.UpdatesStableCoinProposal:
			return k.UpdatesStableCoinProposal(ctx, c)
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s proposal content type: %T", types.ModuleName, c)
		}
	}
}
