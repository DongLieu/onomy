package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/consensys/gnark/backend/groth16"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/onomyprotocol/onomy/x/gnark/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the Msg service.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// CreateCircuit stores the verifying key and circuit artifacts on-chain after validation.
func (m msgServer) CreateCircuit(goCtx context.Context, msg *types.MsgCreateCircuit) (*types.MsgCreateCircuitResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if _, exists := m.GetCircuit(ctx, msg.CircuitId); exists {
		return nil, errorsmod.Wrapf(types.ErrCircuitExists, "circuit %s already exists", msg.CircuitId)
	}

	curveID, err := types.ParseCurveID(msg.CurveId)
	if err != nil {
		return nil, err
	}

	if _, err := types.DecodeVerifyingKey(curveID, msg.VerifyingKey); err != nil {
		return nil, errorsmod.Wrap(err, "failed to decode verifying key")
	}

	curveName, err := types.CurveName(curveID)
	if err != nil {
		return nil, err
	}

	circuit := types.Circuit{
		Id:           msg.CircuitId,
		Creator:      msg.Creator,
		CurveId:      curveName,
		VerifyingKey: msg.VerifyingKey,
		Circuit:      msg.Circuit,
		Height:       uint64(ctx.BlockHeight()),
	}

	m.SetCircuit(ctx, circuit)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCircuitCreated,
			sdk.NewAttribute(types.AttributeKeyCircuitID, msg.CircuitId),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeKeyCurveID, curveName),
		),
	)

	return &types.MsgCreateCircuitResponse{}, nil
}

// VerifyProof checks the provided groth16 proof against the stored verifying key.
func (m msgServer) VerifyProof(goCtx context.Context, msg *types.MsgVerifyProof) (*types.MsgVerifyProofResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	circuit, exists := m.GetCircuit(ctx, msg.CircuitId)
	if !exists {
		return nil, errorsmod.Wrapf(types.ErrCircuitNotFound, "circuit %s", msg.CircuitId)
	}

	curveID, err := types.ParseCurveID(circuit.CurveId)
	if err != nil {
		return nil, err
	}

	vk, err := types.DecodeVerifyingKey(curveID, circuit.VerifyingKey)
	if err != nil {
		return nil, errorsmod.Wrap(err, "stored verifying key cannot be decoded")
	}

	proof, err := types.DecodeProof(curveID, msg.Proof)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid proof")
	}

	publicWitness, err := types.DecodePublicWitness(curveID, msg.PublicWitness)
	if err != nil {
		return nil, errorsmod.Wrap(err, "invalid public witness")
	}

	publicOnly, err := publicWitness.Public()
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to extract public witness")
	}

	if err := groth16.Verify(proof, vk, publicOnly); err != nil {
		return nil, errorsmod.Wrap(types.ErrVerificationFailed, err.Error())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProofVerified,
			sdk.NewAttribute(types.AttributeKeyCircuitID, msg.CircuitId),
			sdk.NewAttribute(types.AttributeKeyCreator, msg.Creator),
		),
	)

	return &types.MsgVerifyProofResponse{Valid: true}, nil
}
