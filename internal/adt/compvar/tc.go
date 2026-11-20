package compvar

import (
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/symbol"
)

func ConvertRecToRef(rec LinearRec) compsem.SemRef {
	return rec.CompRef
}

func ConvertRecsToRefs(recs []LinearRec) []compsem.SemRef {
	if recs == nil {
		return nil
	}
	refs := make([]compsem.SemRef, len(recs))
	for i, rec := range recs {
		refs[i] = ConvertRecToRef(rec)
	}
	return refs
}

func ConvertRecsToRecMap[T VarRec](recs []T) map[symbol.ADT]T {
	recMap := make(map[symbol.ADT]T, len(recs))
	for _, rec := range recs {
		recMap[rec.GetChnlPH()] = rec
	}
	return recMap
}
