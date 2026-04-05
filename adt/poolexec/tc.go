package poolexec

import (
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/valkey"
)

func ConvertRecToRef(rec ExecRec) implsem.SemRef {
	return rec.ImplRef
}

func ExtractExpVKs[T implvar.VarRec](vars map[symbol.ADT]T) map[symbol.ADT]valkey.ADT {
	vks := make(map[symbol.ADT]valkey.ADT, len(vars))
	for _, v := range vars {
		if v.GetChnlPH() == symbol.Zero {
			continue
		}
		vks[v.GetChnlPH()] = v.GetExpVK()
	}
	return vks
}

func DataToExecLiabSnap(dto execLiabSnapDS) (ExecLiabSnap, error) {
	ref, err := implsem.DataToRef(dto.ImplRef)
	if err != nil {
		return ExecLiabSnap{}, err
	}
	mode := implvar.Mode(dto.LiabMode)
	switch mode {
	case implvar.StructMode:
		rec, err := implvar.DataToStructRec(*dto.StructVar)
		if err != nil {
			return ExecLiabSnap{}, err
		}
		return ExecLiabSnap{ref, rec}, nil
	case implvar.LinearMode:
		rec, err := implvar.DataToLinearRec(*dto.LinearVar)
		if err != nil {
			return ExecLiabSnap{}, err
		}
		return ExecLiabSnap{ref, rec}, nil
	default:
		panic(implvar.ErrUnexpectedMode(mode))
	}
}
