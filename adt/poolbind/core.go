package poolbind

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
)

type BindSpec struct {
	// channel placeholder (aka variable name)
	ChnlPH symbol.ADT
	// xact qualified name (aka variable type)
	XactQN uniqsym.ADT
}

type BindRec struct {
	ChnlPH symbol.ADT
	// xact definition id
	DefID identity.ADT
}
