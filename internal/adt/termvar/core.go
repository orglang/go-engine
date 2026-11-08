package termvar

import (
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/adt/valkey"
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
