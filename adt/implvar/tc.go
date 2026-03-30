package implvar

import (
	"database/sql"

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

func ConvertSideToNullInt(val side) sql.NullInt16 {
	return sql.NullInt16{Int16: int16(val), Valid: true}
}

func ConvertSideFromNullInt(val sql.NullInt16) side {
	if val.Valid {
		return side(val.Int16)
	}
	return unkSide
}

func ConvertModeToNullInt(val Mode) sql.NullInt16 {
	return sql.NullInt16{Int16: int16(val), Valid: true}
}

func ConvertModeFromNullInt(val sql.NullInt16) Mode {
	if val.Valid {
		return Mode(val.Int16)
	}
	return unkMode
}
