package poolctx

import (
	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/xactexp"
)

type CtxQry struct {
	DescRef descsem.SemRef
}

type CtxSnap struct {
	DescRef    descsem.SemRef
	StructVars map[symbol.ADT]xactexp.ExpRec
	LinearVars map[symbol.ADT]xactexp.ExpRec
}
