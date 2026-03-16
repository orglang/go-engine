package implvar

import (
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/symbol"
)

func ConvertRecToRef(rec VarRec) implsem.SemRef {
	return rec.ImplRef
}

func ConvertRecsToRecMap(recs []VarRec) map[symbol.ADT]VarRec {
	recMap := make(map[symbol.ADT]VarRec, len(recs))
	for _, rec := range recs {
		recMap[rec.ChnlPH] = rec
	}
	return recMap
}

func DataFromRecMap(recMap map[symbol.ADT]VarRec) []VarRecDS {
	dtos := make([]VarRecDS, 0, len(recMap))
	for _, rec := range recMap {
		dtos = append(dtos, DataFromVarRec(rec))
	}
	return dtos
}

func DataToRecMap(dtos []VarRecDS) (map[symbol.ADT]VarRec, error) {
	recMap := make(map[symbol.ADT]VarRec, len(dtos))
	for _, dto := range dtos {
		rec, err := DataToVarRec(dto)
		if err != nil {
			return nil, err
		}
		recMap[rec.ChnlPH] = rec
	}
	return recMap, nil
}
