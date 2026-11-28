package adapter

import (
	execrootadapter "orglang/go-engine/internal/adt/execroot/adapter"
	"orglang/go-engine/internal/adt/execvar"
	execvaradapter "orglang/go-engine/internal/adt/execvar/adapter"

	execeval "orglang/go-engine/internal/pool/execeval/core"
)

func DataToEvalSnap1(dto EvalSnap1) (execeval.EvalSnap1, error) {
	ref, err := execrootadapter.DataToRef(dto.RootRef)
	if err != nil {
		return execeval.EvalSnap1{}, err
	}
	mode := execvar.Mode(dto.LiabMode)
	switch mode {
	case execvar.StructMode:
		rec, err := execvaradapter.DataToStructRec(*dto.StructVar)
		if err != nil {
			return execeval.EvalSnap1{}, err
		}
		return execeval.EvalSnap1{RootRef: ref, LiabVar: rec}, nil
	case execvar.LinearMode:
		rec, err := execvaradapter.DataToLinearRec(*dto.LinearVar)
		if err != nil {
			return execeval.EvalSnap1{}, err
		}
		return execeval.EvalSnap1{RootRef: ref, LiabVar: rec}, nil
	default:
		panic(execvar.ErrUnexpectedMode(mode))
	}
}
