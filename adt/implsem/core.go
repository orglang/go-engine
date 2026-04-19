package implsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/uniqsym"
)

type SemBind struct {
	ImplQN uniqsym.ADT
	ImplID identity.ADT
}
