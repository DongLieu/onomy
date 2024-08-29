package cli

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/onomyprotocol/onomy/x/psm/types"
)

type AddRequest struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title       string       `json:"title" yaml:"title"`
	Description string       `json:"description" yaml:"description"`
	Denom       string       `json:"denom" yaml:"denom"`
	LimitTotal  sdk.Int      `json:"limit_total" yaml:"limit_total"`
	Price       sdk.Dec      `json:"price" yaml:"price"`
	FeeIn       sdk.Dec      `json:"fee_in" yaml:"fee_in"`
	FeeOut      sdk.Dec      `json:"fee_out" yaml:"fee_out"`
	Deposit     sdk.Coins    `json:"deposit" yaml:"deposit"`
}

// CancelRequest defines a proposal to cancel a current plan.
type UpdatesRequest struct {
	BaseReq          rest.BaseReq `json:"base_req" yaml:"base_req"`
	Title            string       `json:"title" yaml:"title"`
	Description      string       `json:"description" yaml:"description"`
	Denom            string       `json:"denom" yaml:"denom"`
	LimitTotalUpdate sdk.Int      `json:"limit_total_update" yaml:"limit_total_update"`
	Price            sdk.Dec      `json:"price" yaml:"price"`
	FeeIn            sdk.Dec      `json:"fee_in" yaml:"fee_in"`
	FeeOut           sdk.Dec      `json:"fee_out" yaml:"fee_out"`
	Deposit          sdk.Coins    `json:"deposit" yaml:"deposit"`
}

func ProposalRESTAddHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "add_stablecoin",
		Handler:  addProposalHandlerFn(clientCtx),
	}
}

func ProposalRESTUpdateHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "update_stablecoin",
		Handler:  updateProposalHandlerFn(clientCtx),
	}
}

func addProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AddRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// content := proposal.NewParameterChangeProposal(req.Title, req.Description, req.Changes.ToParamChanges())
		content := types.NewAddStableCoinProposal(req.Title, req.Description, req.Denom, req.LimitTotal, req.Price, req.FeeIn, req.FeeOut)
		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		msg, err := govtypes.NewMsgSubmitProposal(&content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

func updateProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdatesRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// content := proposal.NewParameterChangeProposal(req.Title, req.Description, req.Changes.ToParamChanges())
		content := types.NewUpdatesStableCoinProposal(req.Title, req.Description, req.Denom, req.LimitTotalUpdate, req.Price, req.FeeIn, req.FeeOut)
		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		msg, err := govtypes.NewMsgSubmitProposal(&content, req.Deposit, fromAddr)
		if rest.CheckBadRequestError(w, err) {
			return
		}
		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
			return
		}

		tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}
