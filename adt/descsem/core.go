package descsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
	"orglang/go-engine/adt/uniqsym"
)

type SemRef struct {
	DescID identity.ADT
	DescRN seqnum.ADT
}

func NewRef() SemRef {
	return SemRef{identity.New(), seqnum.New()}
}

type SemBind struct {
	DescQN uniqsym.ADT
	DescID identity.ADT
}

type SemRec struct {
	DescRef SemRef
	DescQN  uniqsym.ADT
	Kind    semKind
}

type semKind uint16

const (
	unkKind semKind = iota
	Xact
	Pool
	Type
	Proc
)
