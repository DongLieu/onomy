package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = (*MsgCreateCircuit)(nil)
var _ sdk.Msg = (*MsgVerifyProof)(nil)

// NewMsgCreateCircuit creates a new MsgCreateCircuit instance.
func NewMsgCreateCircuit(creator, id, curveID string, vk, circuit []byte) *MsgCreateCircuit {
	return &MsgCreateCircuit{
		Creator:      creator,
		CircuitId:    id,
		CurveId:      curveID,
		VerifyingKey: vk,
		Circuit:      circuit,
	}
}

// Route implements sdk.Msg.
func (msg MsgCreateCircuit) Route() string { return RouterKey }

// Type implements sdk.Msg.
func (msg MsgCreateCircuit) Type() string { return "CreateCircuit" }

// ValidateBasic implements sdk.Msg.
func (msg MsgCreateCircuit) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return err
	}

	if strings.TrimSpace(msg.CircuitId) == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("circuit_id is required")
	}

	if strings.TrimSpace(msg.CurveId) == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("curve_id is required")
	}

	if len(msg.VerifyingKey) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("verifying_key is required")
	}

	if len(msg.Circuit) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("circuit bytes are required")
	}

	return nil
}

// GetSigners implements sdk.Msg.
func (msg MsgCreateCircuit) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// NewMsgVerifyProof creates a new MsgVerifyProof instance.
func NewMsgVerifyProof(creator, id string, proof, witness []byte) *MsgVerifyProof {
	return &MsgVerifyProof{
		Creator:       creator,
		CircuitId:     id,
		Proof:         proof,
		PublicWitness: witness,
	}
}

// Route implements sdk.Msg.
func (msg MsgVerifyProof) Route() string { return RouterKey }

// Type implements sdk.Msg.
func (msg MsgVerifyProof) Type() string { return "VerifyProof" }

// ValidateBasic implements sdk.Msg.
func (msg MsgVerifyProof) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return err
	}

	if strings.TrimSpace(msg.CircuitId) == "" {
		return sdkerrors.ErrInvalidRequest.Wrap("circuit_id is required")
	}

	if len(msg.Proof) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("proof bytes are required")
	}

	if len(msg.PublicWitness) == 0 {
		return sdkerrors.ErrInvalidRequest.Wrap("public_witness bytes are required")
	}

	return nil
}

// GetSigners implements sdk.Msg.
func (msg MsgVerifyProof) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
