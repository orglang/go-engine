package poolcfg

import (
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/symbol"
)

type CfgSpec struct {
}

type CfgSnap struct {
	SharedVars   map[symbol.ADT]implvar.VarRec
	ProviderVars map[symbol.ADT]implvar.VarRec
	ClientVars   map[symbol.ADT]implvar.VarRec
}
