package implsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/uniqsym"
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
	CommKind
)
