package core

import (
	"orglang/go-engine/internal/adt/contkey"
)

// ExpContentKey hashes the process expression model used by the runtime.
func ExpContentKey(spec ExpSpec) (contkey.ADT, []byte, error) {
	if err := validateExprSpec(spec); err != nil {
		return contkey.Zero, nil, err
	}
	return contkey.HashJSON(spec)
}
