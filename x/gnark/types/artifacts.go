package types

import (
	"bytes"
	"io"

	"cosmossdk.io/errors"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

// DecodeVerifyingKey attempts to deserialize the verifying key bytes for the provided curve.
func DecodeVerifyingKey(curve ecc.ID, bz []byte) (groth16.VerifyingKey, error) {
	vk := groth16.NewVerifyingKey(curve)
	if _, err := vk.ReadFrom(bytes.NewReader(bz)); err != nil {
		return nil, errors.Wrap(ErrInvalidArtifact, err.Error())
	}
	return vk, nil
}

// DecodeProof attempts to deserialize a groth16 proof for the provided curve.
func DecodeProof(curve ecc.ID, bz []byte) (groth16.Proof, error) {
	proof := groth16.NewProof(curve)
	if _, err := proof.ReadFrom(bytes.NewReader(bz)); err != nil {
		return nil, errors.Wrap(ErrInvalidArtifact, err.Error())
	}
	return proof, nil
}

// DecodePublicWitness parses the witness bytes for the provided curve.
func DecodePublicWitness(curve ecc.ID, bz []byte) (witness.Witness, error) {
	w, err := witness.New(curve.ScalarField())
	if err != nil {
		return nil, errors.Wrap(ErrInvalidArtifact, err.Error())
	}
	if _, err := w.ReadFrom(bytes.NewReader(bz)); err != nil {
		return nil, errors.Wrap(ErrInvalidArtifact, err.Error())
	}
	return w, nil
}

// ValidateCircuitBytes ensures the provided R1CS bytes can be decoded for the given curve.
func ValidateCircuitBytes(curve ecc.ID, bz []byte) error {
	cs := groth16.NewCS(curve)
	reader, ok := cs.(io.ReaderFrom)
	if !ok {
		return errors.Wrapf(ErrInvalidArtifact, "curve %d does not support decoding", curve)
	}
	if _, err := reader.ReadFrom(bytes.NewReader(bz)); err != nil {
		return errors.Wrap(ErrInvalidArtifact, err.Error())
	}
	return nil
}
