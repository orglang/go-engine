package exp

import (
	"orglang/go-engine/internal/adt/contkey"
)

// ExpContentKey hashes the execution-oriented expression model used by the runtime.
func ExpContentKey(spec ExpSpec) (contkey.ADT, []byte, error) {
	if err := validateExpSpec(spec); err != nil {
		return contkey.Zero, nil, err
	}
	return contkey.HashJSON(spec)
}
