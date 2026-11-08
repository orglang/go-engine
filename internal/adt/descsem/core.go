package descsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/uniqsym"
)

type SemRec struct {
	DescQN uniqsym.ADT
	DescID identity.ADT
	Kind   descKind
}

type descKind int16

const (
	unkKind descKind = iota
	TypeKind
	TermKind
)
