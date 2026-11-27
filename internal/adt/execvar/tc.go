package execvar

import (
	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/symbol"
)

func ConvertRecToRef(rec LinearRec) comproot.RootRef {
	return rec.RootRef
}

func ConvertRecsToRefs(recs []LinearRec) []comproot.RootRef {
	if recs == nil {
		return nil
	}
	refs := make([]comproot.RootRef, len(recs))
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
