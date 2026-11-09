package termvar

import (
	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
)

// human-readable specification of term variable
// человекочитаемая спецификация переменной терма
type VarSpec struct {
	// channel placeholder (aka variable name)
	ChnlPH symbol.ADT
	// type qualified name (aka variable type)
	TypeQN uniqsym.ADT
}

// machine-readable record of term variable
// машиночитаемая запись переменной терма
type VarRec struct {
	TypeRef typesem.SemRef
	ChnlPH  symbol.ADT
	ExpCK   contkey.ADT
}
