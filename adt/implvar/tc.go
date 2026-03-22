package implvar

import (
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/symbol"
)

func ConvertRecToRef(rec LinearRec) implsem.SemRef {
	return rec.ImplRef
}

func ConvertRecsToRecMap(recs []StructRec) map[symbol.ADT]StructRec {
	recMap := make(map[symbol.ADT]StructRec, len(recs))
	for _, rec := range recs {
		recMap[rec.GetChnlPH()] = rec
	}
	return recMap
}

func DataFromVarRec(rec VarRec) VarRecDS {
	switch varRec := rec.(type) {
	case StructRec:
		return DataFromStructRec(varRec)
	case LinearRec:
		return DataFromLinearRec(varRec)
	default:
		panic(ErrUnexpectedRecType(rec))
	}
}

func DataFromStructMap(recMap map[symbol.ADT]StructRec) []VarRecDS {
	dtos := make([]VarRecDS, 0, len(recMap))
	for _, rec := range recMap {
		dtos = append(dtos, DataFromVarRec(rec))
	}
	return dtos
}

func DataToStructMap(dtos []VarRecDS) (map[symbol.ADT]StructRec, error) {
	recMap := make(map[symbol.ADT]StructRec, len(dtos))
	for _, dto := range dtos {
		rec, err := DataToStructRec(dto)
		if err != nil {
			return nil, err
		}
		recMap[rec.ChnlPH] = rec
	}
	return recMap, nil
}

func DataFromLinearMap(recMap map[symbol.ADT]LinearRec) []VarRecDS {
	dtos := make([]VarRecDS, 0, len(recMap))
	for _, rec := range recMap {
		dtos = append(dtos, DataFromVarRec(rec))
	}
	return dtos
}

func DataToLinearMap(dtos []VarRecDS) (map[symbol.ADT]LinearRec, error) {
	recMap := make(map[symbol.ADT]LinearRec, len(dtos))
	for _, dto := range dtos {
		rec, err := DataToLinearRec(dto)
		if err != nil {
			return nil, err
		}
		recMap[rec.ChnlPH] = rec
	}
	return recMap, nil
}
