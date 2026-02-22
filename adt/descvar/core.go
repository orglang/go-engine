package descvar

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
)

// human-readable specification of description variable
// человекочитаемая спецификация переменной описания
type VarSpec struct {
	// channel placeholder (aka variable name)
	ChnlPH symbol.ADT
	// description qualified name (aka variable type)
	DescQN uniqsym.ADT
}

// machine-readable record of description variable
// машиночитаемая запись переменной описания
type VarRec struct {
	ChnlPH symbol.ADT
	DescID identity.ADT
}
