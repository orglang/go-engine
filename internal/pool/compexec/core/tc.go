package core

import (
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/compvar"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/valkey"
)

func ConvertRecToRef(rec ExecSnap2) compsem.SemRef {
	return rec.CompRef
}

func ExtractExpVKs[T compvar.VarRec](vars []T) map[symbol.ADT]valkey.ADT {
	vks := make(map[symbol.ADT]valkey.ADT, len(vars))
	for _, v := range vars {
		if v.GetChnlPH() == symbol.Zero {
			continue
		}
		vks[v.GetChnlPH()] = v.GetExpVK()
	}
	return vks
}
