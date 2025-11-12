package types

import "fmt"

// DefaultGenesis returns the default module genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:   DefaultGenesisParams(),
		Circuits: []Circuit{},
	}
}

// Validate performs basic genesis validation returning an error upon failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	seen := make(map[string]struct{})
	for _, circuit := range gs.Circuits {
		if circuit.Id == "" {
			return fmt.Errorf("circuit id must not be empty")
		}
		if _, exists := seen[circuit.Id]; exists {
			return fmt.Errorf("duplicate circuit id %q", circuit.Id)
		}
		seen[circuit.Id] = struct{}{}

		if circuit.CurveId == "" {
			return fmt.Errorf("curve_id must be provided for circuit %q", circuit.Id)
		}
		if _, err := ParseCurveID(circuit.CurveId); err != nil {
			return fmt.Errorf("invalid curve for circuit %q: %w", circuit.Id, err)
		}

		if len(circuit.VerifyingKey) == 0 {
			return fmt.Errorf("verifying key missing for circuit %q", circuit.Id)
		}
		if len(circuit.Circuit) == 0 {
			return fmt.Errorf("circuit bytes missing for circuit %q", circuit.Id)
		}
	}

	return nil
}
