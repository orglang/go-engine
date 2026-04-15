package semtype

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/seqnum"
	"orglang/go-engine/adt/uniqsym"
)

type TypeRef struct {
	TypeID identity.ADT
	TypeRN seqnum.ADT
}

func NewRef() TypeRef {
	return TypeRef{identity.New(), seqnum.New()}
}

type SemBind struct {
	DescQN uniqsym.ADT
	DescID identity.ADT
}

type SemRec struct {
	DescRef TypeRef
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
