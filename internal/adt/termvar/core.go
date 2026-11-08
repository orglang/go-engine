package termvar

import (
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/typesem"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
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
	ExpVK   valkey.ADT
}
