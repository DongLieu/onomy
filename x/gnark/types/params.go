package types

// DefaultGenesisParams returns default module parameters.
func DefaultGenesisParams() Params {
	return Params{}
}

// Validate performs basic parameter validation.
func (p Params) Validate() error {
	return nil
}
