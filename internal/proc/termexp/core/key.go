package core

import "orglang/go-engine/internal/adt/contkey"

func ContentKey(spec ExpSpec) (contkey.ADT, []byte, error) {
	if err := validateDefinitionSpec(spec); err != nil {
		return contkey.Zero, nil, err
	}
	return contkey.HashJSON(spec)
}
