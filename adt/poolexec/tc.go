package poolexec

import (
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
)

func ConvertRecToRef(rec ExecRec) implsem.SemRef {
	return rec.ImplRef
}

func DataToExecLiabSnap(dto execSnapDS) (ExecSnap, error) {
	ref, err := implsem.DataToRef(dto.ImplRef)
	if err != nil {
		return ExecSnap{}, err
	}
	mode := implvar.Mode(dto.LiabMode)
	switch mode {
	case implvar.StructMode:
		rec, err := implvar.DataToStructRec(*dto.StructVar)
		if err != nil {
			return ExecSnap{}, err
		}
		return ExecSnap{ref, rec}, nil
	case implvar.LinearMode:
		rec, err := implvar.DataToLinearRec(*dto.LinearVar)
		if err != nil {
			return ExecSnap{}, err
		}
		return ExecSnap{ref, rec}, nil
	default:
		panic(implvar.ErrUnexpectedMode(mode))
	}
}
