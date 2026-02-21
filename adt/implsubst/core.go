package implsubst

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
)

type SubstSpec struct {
	ChnlPH symbol.ADT
	ImplQN uniqsym.ADT
}

type SubstRec struct {
	ChnlPH symbol.ADT
	ImplID identity.ADT
}
