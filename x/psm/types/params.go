package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	DefaultLimitTotal           = sdk.NewInt(100_000_000)
	DefaultAcceptablePriceRatio = sdk.MustNewDecFromStr("0.001")
)

// Parameter store keys.
var (
	// KeyPoolRate is byte key for KeyPoolRate param.
	KeyLimmitTotal          = []byte("LimitTotal") //nolint:gochecknoglobals // cosmos-sdk style
	KeyAcceptablePriceRatio = []byte("AcceptablePriceRatio")
)

// NewParams creates a new Params instance.
func NewParams(
	limitTotal sdk.Int,
	AcceptablePriceRatio sdk.Dec,
) Params {
	return Params{
		LimitTotal:           limitTotal,
		AcceptablePriceRatio: AcceptablePriceRatio,
	}
}

func DefaultParams() Params {
	return NewParams(
		DefaultLimitTotal, DefaultAcceptablePriceRatio,
	)
}

// ParamSetPairs get the params.ParamSet.
func (m *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLimmitTotal, &m.LimitTotal, validateLimitTotal),
		paramtypes.NewParamSetPair(KeyAcceptablePriceRatio, &m.AcceptablePriceRatio, validateAcceptablePriceRatio),
	}
}

// Validate validates the set of params.
func (m Params) Validate() error {
	if err := validateLimitTotal(m.LimitTotal); err != nil {
		return err
	}
	if m.AcceptablePriceRatio.LTE(sdk.ZeroDec()) {
		return fmt.Errorf("AcceptablePriceRatio must be positive")
	}
	return nil
}

func validateAcceptablePriceRatio(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() || v.IsNegative() {
		return fmt.Errorf("AcceptablePriceRatio cannot be negative or nil: %s", v)
	}

	return nil
}

func validateLimitTotal(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() || v.IsNegative() {
		return fmt.Errorf("total limit rate cannot be negative or nil: %s", v)
	}

	return nil
}

// ParamTable for staking module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}
