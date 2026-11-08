package descsem

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
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
