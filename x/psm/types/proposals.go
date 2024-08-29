package types

import (
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeAddStableCoinProposal     string = "AddStableCoinProposal"
	ProposalTypeUpdatesStableCoinProposal string = "UpdatesStableCoinProposal"
)

var (
	_ govtypes.Content = &AddStableCoinProposal{}
	_ govtypes.Content = &UpdatesStableCoinProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeAddStableCoinProposal)
	govtypes.RegisterProposalTypeCodec(&AddStableCoinProposal{}, fmt.Sprintf("%s/%s", ModuleName, ProposalTypeAddStableCoinProposal))

	govtypes.RegisterProposalType(ProposalTypeUpdatesStableCoinProposal)
	govtypes.RegisterProposalTypeCodec(&UpdatesStableCoinProposal{}, fmt.Sprintf("%s/%s", ModuleName, ProposalTypeUpdatesStableCoinProposal))
}

func NewAddStableCoinProposal(title, description, denom string, limitTotal sdk.Int, price, feeIn, feeOut sdk.Dec) AddStableCoinProposal {
	return AddStableCoinProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
		LimitTotal:  limitTotal,
		Price:       price,
		FeeIn:       feeIn,
		FeeOut:      feeOut,
	}
}

func NewUpdatesStableCoinProposal(title, description, denom string, updateLimitTotal sdk.Int, price, feeIn, feeOut sdk.Dec) UpdatesStableCoinProposal {
	return UpdatesStableCoinProposal{
		Title:             title,
		Description:       description,
		Denom:             denom,
		UpdatesLimitTotal: updateLimitTotal,
		Price:             price,
		FeeIn:             feeIn,
		FeeOut:            feeOut,
	}
}

func (a *AddStableCoinProposal) ProposalRoute() string { return RouterKey }

func (a *AddStableCoinProposal) ProposalType() string {
	return ProposalTypeAddStableCoinProposal
}

func (a *AddStableCoinProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(a)
	if err != nil {
		return err
	}

	if a.Denom == "" {
		return sdkerrors.Wrap(ErrInvalidAddStableCoinProposal, "empty denom")
	}
	if a.LimitTotal.LT(sdk.ZeroInt()) {
		return sdkerrors.Wrap(ErrInvalidAddStableCoinProposal, "less than zero")
	}
	return nil
}

// func (a AddStableCoinProposal)
// func (a AddStableCoinProposal)

func (u *UpdatesStableCoinProposal) ProposalRoute() string { return RouterKey }

func (u *UpdatesStableCoinProposal) ProposalType() string {
	return ProposalTypeUpdatesStableCoinProposal
}

func (u *UpdatesStableCoinProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(u)
	if err != nil {
		return err
	}

	if u.Denom == "" {
		return sdkerrors.Wrap(ErrInvalidAddStableCoinProposal, "empty denom")
	}
	if u.UpdatesLimitTotal.LT(sdk.ZeroInt()) {
		return sdkerrors.Wrap(ErrInvalidAddStableCoinProposal, "less than zero")
	}
	return nil
}
