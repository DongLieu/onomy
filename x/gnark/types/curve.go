package types

import (
	"strings"

	"cosmossdk.io/errors"
	"github.com/consensys/gnark-crypto/ecc"
)

var curveNameToID = map[string]ecc.ID{
	"BN254":     ecc.BN254,
	"BLS12_377": ecc.BLS12_377,
	"BLS12_381": ecc.BLS12_381,
	"BW6_761":   ecc.BW6_761,
	"BLS24_315": ecc.BLS24_315,
	"BLS24_317": ecc.BLS24_317,
	"BW6_633":   ecc.BW6_633,
}

// ParseCurveID converts a textual curve identifier into the gnark ecc.ID value.
func ParseCurveID(value string) (ecc.ID, error) {
	replacer := strings.NewReplacer("-", "_", " ", "")
	name := replacer.Replace(strings.ToUpper(strings.TrimSpace(value)))
	if id, ok := curveNameToID[name]; ok {
		return id, nil
	}
	return 0, errors.Wrapf(ErrUnsupportedCurve, "curve %q", value)
}

// CurveName returns the canonical curve name for a given ecc.ID.
func CurveName(id ecc.ID) (string, error) {
	for name, curveID := range curveNameToID {
		if curveID == id {
			return name, nil
		}
	}
	return "", errors.Wrapf(ErrUnsupportedCurve, "curve id %d", id)
}
