package implsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/uniqsym"
)

type SemRef struct {
	ImplID identity.ADT
	ImplRN revnum.ADT
}

func NewRef() SemRef {
	return SemRef{identity.New(), revnum.New()}
}

type SemBind struct {
	ImplQN uniqsym.ADT
	ImplID identity.ADT
}

type SemRec struct {
	Ref  SemRef
	Bind SemBind
	Kind semKind
}

type semKind uint8

const (
	unknown semKind = iota
	Pool
	Proc
)
