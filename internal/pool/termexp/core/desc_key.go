package core

import (
	"orglang/go-engine/internal/adt/contkey"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

func DescContentKey(spec DescSpec) (contkey.ADT, []byte, error) { return termdef.ExpContentKey(spec) }
