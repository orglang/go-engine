package compvar

import (
	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/symbol"
)

func ConvertRecToRef(rec LinearRec) comproot.SemRef {
	return rec.CompRef
}

func ConvertRecsToRefs(recs []LinearRec) []comproot.SemRef {
	if recs == nil {
		return nil
	}
	refs := make([]comproot.SemRef, len(recs))
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
