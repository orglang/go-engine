package adapter

import (
	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/compvar"

	compexec "orglang/go-engine/pool/compexec/core"
)

func DataToExecSnap1(dto compexec.ExecSnap1DS) (compexec.ExecSnap1, error) {
	ref, err := compsem.DataToRef(dto.CompRef)
	if err != nil {
		return compexec.ExecSnap1{}, err
	}
	mode := compvar.Mode(dto.LiabMode)
	switch mode {
	case compvar.StructMode:
		rec, err := compvar.DataToStructRec(*dto.StructVar)
		if err != nil {
			return compexec.ExecSnap1{}, err
		}
		return compexec.ExecSnap1{ref, rec}, nil
	case compvar.LinearMode:
		rec, err := compvar.DataToLinearRec(*dto.LinearVar)
		if err != nil {
			return compexec.ExecSnap1{}, err
		}
		return compexec.ExecSnap1{ref, rec}, nil
	default:
		panic(compvar.ErrUnexpectedMode(mode))
	}
}
