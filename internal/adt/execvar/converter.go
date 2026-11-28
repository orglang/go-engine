package execvar

import (
	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/symbol"
)

func ConvertRecToRef(rec LinearRec) execroot.RootRef {
	return rec.RootRef
}

func ConvertRecsToRefs(recs []LinearRec) []execroot.RootRef {
	if recs == nil {
		return nil
	}
	refs := make([]execroot.RootRef, len(recs))
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
