package adapter

import (
	comprootadapter "orglang/go-engine/internal/adt/comproot/adapter"
	"orglang/go-engine/internal/adt/compvar"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"

	compexec "orglang/go-engine/internal/pool/compexec/core"
)

func DataToExecSnap1(dto EvalSnap1) (compexec.EvalSnap1, error) {
	ref, err := comprootadapter.DataToRef(dto.CompRef)
	if err != nil {
		return compexec.EvalSnap1{}, err
	}
	mode := compvar.Mode(dto.LiabMode)
	switch mode {
	case compvar.StructMode:
		rec, err := compvaradapter.DataToStructRec(*dto.StructVar)
		if err != nil {
			return compexec.EvalSnap1{}, err
		}
		return compexec.EvalSnap1{CompRef: ref, LiabVar: rec}, nil
	case compvar.LinearMode:
		rec, err := compvaradapter.DataToLinearRec(*dto.LinearVar)
		if err != nil {
			return compexec.EvalSnap1{}, err
		}
		return compexec.EvalSnap1{CompRef: ref, LiabVar: rec}, nil
	default:
		panic(compvar.ErrUnexpectedMode(mode))
	}
}
