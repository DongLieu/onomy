package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrCircuitExists      = errors.Register(ModuleName, 1, "circuit already exists")
	ErrCircuitNotFound    = errors.Register(ModuleName, 2, "circuit not found")
	ErrUnsupportedCurve   = errors.Register(ModuleName, 3, "unsupported curve id")
	ErrInvalidArtifact    = errors.Register(ModuleName, 4, "invalid artifact")
	ErrVerificationFailed = errors.Register(ModuleName, 5, "verification failed")
)
