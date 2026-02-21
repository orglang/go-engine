package implsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
)

type SemRef struct {
	ImplID identity.ADT
	ImplRN revnum.ADT
}

func NewRef() SemRef {
	return SemRef{identity.New(), revnum.New()}
}

type SemRec struct {
	Ref  SemRef
	Kind semKind
}

type semKind uint8

const (
	unknown semKind = iota
	Pool
	Proc
)
