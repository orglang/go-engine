package adapter

import (
	"database/sql"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/execvar"
	contkeyadapter "orglang/go-engine/internal/adt/contkey/adapter"
	identityadapter "orglang/go-engine/internal/adt/identity/adapter"
	seqnumadapter "orglang/go-engine/internal/adt/seqnum/adapter"
	"orglang/go-engine/internal/adt/symbol"
	symboladapter "orglang/go-engine/internal/adt/symbol/adapter"
)

func DataFromVarRec(rec execvar.VarRec) VarRec {
	switch value := rec.(type) {
	case execvar.StructRec:
		return DataFromStructRec(value)
	case execvar.LinearRec:
		return DataFromLinearRec(value)
	default:
		panic(execvar.ErrUnexpectedRecType(rec))
	}
}

func DataFromStructRec(rec execvar.StructRec) VarRec {
	return VarRec{
		ExecID: identityadapter.ConvertToNullString(rec.ExecRef.ExecID),
		ExecRN: seqnumadapter.ConvertToNullInt(rec.ExecRef.ExecRN),
		ConnID: identityadapter.ConvertToNullString(rec.ConnID),
		ChnlID: identityadapter.ConvertToNullString(rec.ChnlID),
		ChnlPH: symboladapter.ConvertToNullString(rec.ChnlPH),
		ExpCK:  contkeyadapter.ConvertToNullInt(rec.ExpCK),
		ChnlBS: ConvertSideToNullInt(rec.ChnlBS),
	}
}

func DataFromLinearRec(rec execvar.LinearRec) VarRec {
	return VarRec{
		ExecID: identityadapter.ConvertToNullString(rec.ExecRef.ExecID),
		ExecRN: seqnumadapter.ConvertToNullInt(rec.ExecRef.ExecRN),
		ConnID: identityadapter.ConvertToNullString(rec.ConnID),
		ChnlID: identityadapter.ConvertToNullString(rec.ChnlID),
		ChnlPH: symboladapter.ConvertToNullString(rec.ChnlPH),
		ExpCK:  contkeyadapter.ConvertToNullInt(rec.ExpCK),
		ChnlBS: ConvertSideToNullInt(rec.ChnlBS),
	}
}

func DataFromStructMap(recMap map[symbol.ADT]execvar.StructRec) []VarRec {
	dtos := make([]VarRec, 0, len(recMap))
	for _, rec := range recMap {
		dtos = append(dtos, DataFromStructRec(rec))
	}
	return dtos
}

func DataToStructMap(dtos []VarRec) (map[symbol.ADT]execvar.StructRec, error) {
	recMap := make(map[symbol.ADT]execvar.StructRec, len(dtos))
	for _, dto := range dtos {
		rec, err := DataToStructRec(dto)
		if err != nil {
			return nil, err
		}
		recMap[rec.ChnlPH] = rec
	}
	return recMap, nil
}

func DataFromLinearMap(recMap map[symbol.ADT]execvar.LinearRec) []VarRec {
	dtos := make([]VarRec, 0, len(recMap))
	for _, rec := range recMap {
		dtos = append(dtos, DataFromLinearRec(rec))
	}
	return dtos
}

func DataToLinearMap(dtos []VarRec) (map[symbol.ADT]execvar.LinearRec, error) {
	recMap := make(map[symbol.ADT]execvar.LinearRec, len(dtos))
	for _, dto := range dtos {
		rec, err := DataToLinearRec(dto)
		if err != nil {
			return nil, err
		}
		recMap[rec.ChnlPH] = rec
	}
	return recMap, nil
}

func DataToStructRec(dto VarRec) (execvar.StructRec, error) {
	ref, err := dataToExecRef(dto)
	if err != nil {
		return execvar.StructRec{}, err
	}
	connID, err := identityadapter.ConvertFromNullString(dto.ConnID)
	if err != nil {
		return execvar.StructRec{}, err
	}
	chnlID, err := identityadapter.ConvertFromNullString(dto.ChnlID)
	if err != nil {
		return execvar.StructRec{}, err
	}
	chnlPH, err := symboladapter.ConvertFromNullString(dto.ChnlPH)
	if err != nil {
		return execvar.StructRec{}, err
	}
	return execvar.StructRec{
		ExecRef: ref,
		ConnID:  connID,
		ChnlID:  chnlID,
		ChnlPH:  chnlPH,
		ChnlBS:  ConvertSideFromNullInt(dto.ChnlBS),
		ExpCK:   contkeyadapter.ConvertFromNullInt(dto.ExpCK),
	}, nil
}

func DataToLinearRec(dto VarRec) (execvar.LinearRec, error) {
	ref, err := dataToExecRef(dto)
	if err != nil {
		return execvar.LinearRec{}, err
	}
	connID, err := identityadapter.ConvertFromNullString(dto.ConnID)
	if err != nil {
		return execvar.LinearRec{}, err
	}
	chnlID, err := identityadapter.ConvertFromNullString(dto.ChnlID)
	if err != nil {
		return execvar.LinearRec{}, err
	}
	chnlPH, err := symboladapter.ConvertFromNullString(dto.ChnlPH)
	if err != nil {
		return execvar.LinearRec{}, err
	}
	return execvar.LinearRec{
		ExecRef: ref,
		ConnID:  connID,
		ChnlID:  chnlID,
		ChnlPH:  chnlPH,
		ChnlBS:  ConvertSideFromNullInt(dto.ChnlBS),
		ExpCK:   contkeyadapter.ConvertFromNullInt(dto.ExpCK),
	}, nil
}

func DataFromStructRecs(recs []execvar.StructRec) []VarRec {
	if recs == nil {
		return nil
	}
	dtos := make([]VarRec, len(recs))
	for i, rec := range recs {
		dtos[i] = DataFromStructRec(rec)
	}
	return dtos
}

func DataFromLinearRecs(recs []execvar.LinearRec) []VarRec {
	if recs == nil {
		return nil
	}
	dtos := make([]VarRec, len(recs))
	for i, rec := range recs {
		dtos[i] = DataFromLinearRec(rec)
	}
	return dtos
}

func DataToStructRecs(dtos []VarRec) ([]execvar.StructRec, error) {
	if dtos == nil {
		return nil, nil
	}
	recs := make([]execvar.StructRec, len(dtos))
	for i, dto := range dtos {
		rec, err := DataToStructRec(dto)
		if err != nil {
			return nil, err
		}
		recs[i] = rec
	}
	return recs, nil
}

func DataToLinearRecs(dtos []VarRec) ([]execvar.LinearRec, error) {
	if dtos == nil {
		return nil, nil
	}
	recs := make([]execvar.LinearRec, len(dtos))
	for i, dto := range dtos {
		rec, err := DataToLinearRec(dto)
		if err != nil {
			return nil, err
		}
		recs[i] = rec
	}
	return recs, nil
}

func dataToExecRef(dto VarRec) (comproot.SemRef, error) {
	execID, err := identityadapter.ConvertFromNullString(dto.ExecID)
	if err != nil {
		return comproot.SemRef{}, err
	}
	return comproot.SemRef{
		ExecID: execID,
		ExecRN: seqnumadapter.ConvertFromNullInt(dto.ExecRN),
	}, nil
}

func ConvertSideToNullInt(value execvar.Side) sql.NullInt16 {
	return sql.NullInt16{Int16: int16(value), Valid: true}
}

func ConvertSideFromNullInt(value sql.NullInt16) execvar.Side {
	if value.Valid {
		return execvar.Side(value.Int16)
	}
	return execvar.Side(0)
}

func ConvertModeToNullInt(value execvar.Mode) sql.NullInt16 {
	return sql.NullInt16{Int16: int16(value), Valid: true}
}

func ConvertModeFromNullInt(value sql.NullInt16) execvar.Mode {
	if value.Valid {
		return execvar.Mode(value.Int16)
	}
	return execvar.Mode(0)
}
