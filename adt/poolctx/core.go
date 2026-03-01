package poolctx

import (
	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/xactexp"
)

type CtxSpec struct {
	DescRef descsem.SemRef
}

type CtxSnap struct {
	// активы пула
	ClientVars map[symbol.ADT]xactexp.ExpRec
	// обязательства пула
	ProviderVars map[symbol.ADT]xactexp.ExpRec
}
