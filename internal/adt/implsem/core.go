package implsem

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
)

type SemRec struct {
	ImplQN uniqsym.ADT
	ImplID identity.ADT
	Kind   implKind
}

type implKind int16

const (
	unkKind implKind = iota
	CompKind
	ConnKind
)
