package poolvar

import (
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/xactexp"
)

type VarCtx map[symbol.ADT]xactexp.ExpSpec
