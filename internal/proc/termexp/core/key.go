package core

import (
	"orglang/go-engine/internal/adt/contkey"
)

// CompContentKey hashes the execution-oriented expression model used by the runtime.
func CompContentKey(spec CompSpec) (contkey.ADT, []byte, error) {
	if err := validateCompSpec(spec); err != nil {
		return contkey.Zero, nil, err
	}
	return contkey.HashJSON(spec)
}
