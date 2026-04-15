package compexec

import (
	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/valkey"
)

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
