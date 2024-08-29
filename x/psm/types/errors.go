package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidAddStableCoinProposal = sdkerrors.Register(ModuleName, 2, "invalid add stable coin proposal")
)
