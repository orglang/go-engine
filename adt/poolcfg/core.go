package poolcfg

import (
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/symbol"
)

type CfgSpec struct {
}

type CfgSnap struct {
	SharedVars   map[symbol.ADT]implvar.StructRec
	ProviderVars map[symbol.ADT]implvar.LinearRec
	ClientVars   map[symbol.ADT]implvar.LinearRec
}
