package core

import (
	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/execvar"
	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
)

func ConvertRecToRef(rec EvalSnap2) comproot.RootRef {
	return rec.CompRef
}

func ExtractExpCKs[T execvar.VarRec](vars []T) map[symbol.ADT]contkey.ADT {
	cks := make(map[symbol.ADT]contkey.ADT, len(vars))
	for _, v := range vars {
		if v.GetChnlPH() == symbol.Zero {
			continue
		}
		cks[v.GetChnlPH()] = v.GetExpCK()
	}
	return cks
}
