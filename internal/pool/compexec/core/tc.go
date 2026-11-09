package core

import (
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/compvar"
	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
)

func ConvertRecToRef(rec ExecSnap2) compsem.SemRef {
	return rec.CompRef
}

func ExtractExpCKs[T compvar.VarRec](vars []T) map[symbol.ADT]contkey.ADT {
	cks := make(map[symbol.ADT]contkey.ADT, len(vars))
	for _, v := range vars {
		if v.GetChnlPH() == symbol.Zero {
			continue
		}
		cks[v.GetChnlPH()] = v.GetExpCK()
	}
	return cks
}
