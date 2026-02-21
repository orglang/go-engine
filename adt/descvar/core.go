package descvar

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
)

type VarSpec struct {
	// channel placeholder (aka variable name)
	ChnlPH symbol.ADT
	// description qualified name (aka variable type)
	DescQN uniqsym.ADT
}

type VarRec struct {
	ChnlPH symbol.ADT
	DescID identity.ADT
}
