package descsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/uniqsym"
)

type SemRef struct {
	DescID identity.ADT
	DescRN revnum.ADT
}

func NewRef() SemRef {
	return SemRef{identity.New(), revnum.New()}
}

type SemBind struct {
	DescQN uniqsym.ADT
	DescID identity.ADT
}

type SemRec struct {
	Ref  SemRef
	Bind SemBind
	Kind semKind
}

type semKind uint8

const (
	nonDesc semKind = iota
	Xact
	Pool
	Type
	Proc
)
