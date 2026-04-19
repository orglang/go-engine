package descsem

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/uniqsym"
)

type SemBind struct {
	DescQN uniqsym.ADT
	DescID identity.ADT
}
