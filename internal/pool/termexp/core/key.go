package core

import (
	"orglang/go-engine/internal/adt/contkey"
	compexp "orglang/go-engine/internal/pool/compexec/exp"
)

func CompContentKey(spec CompSpec) (contkey.ADT, []byte, error) { return compexp.ExpContentKey(spec) }
