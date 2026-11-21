package adapter

import (
	"database/sql"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/compvar"
	contkeyadapter "orglang/go-engine/internal/adt/contkey/adapter"
	identityadapter "orglang/go-engine/internal/adt/identity/adapter"
	seqnumadapter "orglang/go-engine/internal/adt/seqnum/adapter"
	"orglang/go-engine/internal/adt/symbol"
	symboladapter "orglang/go-engine/internal/adt/symbol/adapter"
)

func DataFromVarRec(rec compvar.VarRec) VarRecDS {
	switch value := rec.(type) {
	case compvar.StructRec:
		return DataFromStructRec(value)
	case compvar.LinearRec:
		return DataFromLinearRec(value)
	default:
		panic(compvar.ErrUnexpectedRecType(rec))
	}
}

func DataFromStructRec(rec compvar.StructRec) VarRecDS {
	return VarRecDS{
		CompID: identityadapter.ConvertToNullString(rec.CompRef.CompID),
		CompRN: seqnumadapter.ConvertToNullInt(rec.CompRef.CompRN),
		ConnID: identityadapter.ConvertToNullString(rec.ConnID),
		ChnlID: identityadapter.ConvertToNullString(rec.ChnlID),
		ChnlPH: symboladapter.ConvertToNullString(rec.ChnlPH),
		ExpCK:  contkeyadapter.ConvertToNullInt(rec.ExpCK),
		ChnlBS: ConvertSideToNullInt(rec.ChnlBS),
	}
}

func DataFromLinearRec(rec compvar.LinearRec) VarRecDS {
	return VarRecDS{
		CompID: identityadapter.ConvertToNullString(rec.CompRef.CompID),
		CompRN: seqnumadapter.ConvertToNullInt(rec.CompRef.CompRN),
		ConnID: identityadapter.ConvertToNullString(rec.ConnID),
		ChnlID: identityadapter.ConvertToNullString(rec.ChnlID),
		ChnlPH: symboladapter.ConvertToNullString(rec.ChnlPH),
		ExpCK:  contkeyadapter.ConvertToNullInt(rec.ExpCK),
		ChnlBS: ConvertSideToNullInt(rec.ChnlBS),
	}
}

func DataFromStructMap(recMap map[symbol.ADT]compvar.StructRec) []VarRecDS {
	dtos := make([]VarRecDS, 0, len(recMap))
	for _, rec := range recMap {
		dtos = append(dtos, DataFromStructRec(rec))
	}
	return dtos
}

func DataToStructMap(dtos []VarRecDS) (map[symbol.ADT]compvar.StructRec, error) {
	recMap := make(map[symbol.ADT]compvar.StructRec, len(dtos))
	for _, dto := range dtos {
		rec, err := DataToStructRec(dto)
		if err != nil {
			return nil, err
		}
		recMap[rec.ChnlPH] = rec
	}
	return recMap, nil
}

func DataFromLinearMap(recMap map[symbol.ADT]compvar.LinearRec) []VarRecDS {
	dtos := make([]VarRecDS, 0, len(recMap))
	for _, rec := range recMap {
		dtos = append(dtos, DataFromLinearRec(rec))
	}
	return dtos
}

func DataToLinearMap(dtos []VarRecDS) (map[symbol.ADT]compvar.LinearRec, error) {
	recMap := make(map[symbol.ADT]compvar.LinearRec, len(dtos))
	for _, dto := range dtos {
		rec, err := DataToLinearRec(dto)
		if err != nil {
			return nil, err
		}
		recMap[rec.ChnlPH] = rec
	}
	return recMap, nil
}

func DataToStructRec(dto VarRecDS) (compvar.StructRec, error) {
	ref, err := dataToCompRef(dto)
	if err != nil {
		return compvar.StructRec{}, err
	}
	connID, err := identityadapter.ConvertFromNullString(dto.ConnID)
	if err != nil {
		return compvar.StructRec{}, err
	}
	chnlID, err := identityadapter.ConvertFromNullString(dto.ChnlID)
	if err != nil {
		return compvar.StructRec{}, err
	}
	chnlPH, err := symboladapter.ConvertFromNullString(dto.ChnlPH)
	if err != nil {
		return compvar.StructRec{}, err
	}
	return compvar.StructRec{
		CompRef: ref,
		ConnID:  connID,
		ChnlID:  chnlID,
		ChnlPH:  chnlPH,
		ChnlBS:  ConvertSideFromNullInt(dto.ChnlBS),
		ExpCK:   contkeyadapter.ConvertFromNullInt(dto.ExpCK),
	}, nil
}

func DataToLinearRec(dto VarRecDS) (compvar.LinearRec, error) {
	ref, err := dataToCompRef(dto)
	if err != nil {
		return compvar.LinearRec{}, err
	}
	connID, err := identityadapter.ConvertFromNullString(dto.ConnID)
	if err != nil {
		return compvar.LinearRec{}, err
	}
	chnlID, err := identityadapter.ConvertFromNullString(dto.ChnlID)
	if err != nil {
		return compvar.LinearRec{}, err
	}
	chnlPH, err := symboladapter.ConvertFromNullString(dto.ChnlPH)
	if err != nil {
		return compvar.LinearRec{}, err
	}
	return compvar.LinearRec{
		CompRef: ref,
		ConnID:  connID,
		ChnlID:  chnlID,
		ChnlPH:  chnlPH,
		ChnlBS:  ConvertSideFromNullInt(dto.ChnlBS),
		ExpCK:   contkeyadapter.ConvertFromNullInt(dto.ExpCK),
	}, nil
}

func DataFromStructRecs(recs []compvar.StructRec) []VarRecDS {
	if recs == nil {
		return nil
	}
	dtos := make([]VarRecDS, len(recs))
	for i, rec := range recs {
		dtos[i] = DataFromStructRec(rec)
	}
	return dtos
}

func DataFromLinearRecs(recs []compvar.LinearRec) []VarRecDS {
	if recs == nil {
		return nil
	}
	dtos := make([]VarRecDS, len(recs))
	for i, rec := range recs {
		dtos[i] = DataFromLinearRec(rec)
	}
	return dtos
}

func DataToStructRecs(dtos []VarRecDS) ([]compvar.StructRec, error) {
	if dtos == nil {
		return nil, nil
	}
	recs := make([]compvar.StructRec, len(dtos))
	for i, dto := range dtos {
		rec, err := DataToStructRec(dto)
		if err != nil {
			return nil, err
		}
		recs[i] = rec
	}
	return recs, nil
}

func DataToLinearRecs(dtos []VarRecDS) ([]compvar.LinearRec, error) {
	if dtos == nil {
		return nil, nil
	}
	recs := make([]compvar.LinearRec, len(dtos))
	for i, dto := range dtos {
		rec, err := DataToLinearRec(dto)
		if err != nil {
			return nil, err
		}
		recs[i] = rec
	}
	return recs, nil
}

func dataToCompRef(dto VarRecDS) (compsem.SemRef, error) {
	compID, err := identityadapter.ConvertFromNullString(dto.CompID)
	if err != nil {
		return compsem.SemRef{}, err
	}
	return compsem.SemRef{
		CompID: compID,
		CompRN: seqnumadapter.ConvertFromNullInt(dto.CompRN),
	}, nil
}

func ConvertSideToNullInt(value compvar.Side) sql.NullInt16 {
	return sql.NullInt16{Int16: int16(value), Valid: true}
}

func ConvertSideFromNullInt(value sql.NullInt16) compvar.Side {
	if value.Valid {
		return compvar.Side(value.Int16)
	}
	return compvar.Side(0)
}

func ConvertModeToNullInt(value compvar.Mode) sql.NullInt16 {
	return sql.NullInt16{Int16: int16(value), Valid: true}
}

func ConvertModeFromNullInt(value sql.NullInt16) compvar.Mode {
	if value.Valid {
		return compvar.Mode(value.Int16)
	}
	return compvar.Mode(0)
}
