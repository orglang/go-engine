package compexec

import (
	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/valkey"
)

func ConvertRecToRef(rec ExecSnap2) compsem.SemRef {
	return rec.CompRef
}

func ExtractExpVKs[T compvar.VarRec](vars []T) map[symbol.ADT]valkey.ADT {
	vks := make(map[symbol.ADT]valkey.ADT, len(vars))
	for _, v := range vars {
		if v.GetChnlPH() == symbol.Zero {
			continue
		}
		vks[v.GetChnlPH()] = v.GetExpVK()
	}
	return vks
}

func DataToExecSnap1(dto execSnap1) (ExecSnap1, error) {
	ref, err := compsem.DataToRef(dto.CompRef)
	if err != nil {
		return ExecSnap1{}, err
	}
	mode := compvar.Mode(dto.LiabMode)
	switch mode {
	case compvar.StructMode:
		rec, err := compvar.DataToStructRec(*dto.StructVar)
		if err != nil {
			return ExecSnap1{}, err
		}
		return ExecSnap1{ref, rec}, nil
	case compvar.LinearMode:
		rec, err := compvar.DataToLinearRec(*dto.LinearVar)
		if err != nil {
			return ExecSnap1{}, err
		}
		return ExecSnap1{ref, rec}, nil
	default:
		panic(compvar.ErrUnexpectedMode(mode))
	}
}
