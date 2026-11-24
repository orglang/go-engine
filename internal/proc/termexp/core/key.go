package core

import (
	"fmt"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
)

// CompContentKey hashes the execution-oriented expression model used by the runtime.
func CompContentKey(spec CompSpec) (contkey.ADT, []byte, error) {
	if err := validateCompSpec(spec); err != nil {
		return contkey.Zero, nil, err
	}
	return contkey.HashJSON(spec)
}
