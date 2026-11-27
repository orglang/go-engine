package adapter

import (
	comprootadapter "orglang/go-engine/internal/adt/comproot/adapter"
	"orglang/go-engine/internal/adt/compvar"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"

	execeval "orglang/go-engine/internal/pool/execeval/core"
)

func DataToEvalSnap1(dto EvalSnap1) (execeval.EvalSnap1, error) {
	ref, err := comprootadapter.DataToRef(dto.CompRef)
	if err != nil {
		return execeval.EvalSnap1{}, err
	}
	mode := compvar.Mode(dto.LiabMode)
	switch mode {
	case compvar.StructMode:
		rec, err := compvaradapter.DataToStructRec(*dto.StructVar)
		if err != nil {
			return execeval.EvalSnap1{}, err
		}
		return execeval.EvalSnap1{CompRef: ref, LiabVar: rec}, nil
	case compvar.LinearMode:
		rec, err := compvaradapter.DataToLinearRec(*dto.LinearVar)
		if err != nil {
			return execeval.EvalSnap1{}, err
		}
		return execeval.EvalSnap1{CompRef: ref, LiabVar: rec}, nil
	default:
		panic(compvar.ErrUnexpectedMode(mode))
	}
}
