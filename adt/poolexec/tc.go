package poolexec

import (
	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/semterm"
)

func ConvertRecToRef(rec ExecRec) semterm.TermRef {
	return rec.CompRef
}

func DataToExecLiabSnap(dto execSnapDS) (ExecSnap, error) {
	ref, err := semterm.DataToRef(dto.ImplRef)
	if err != nil {
		return ExecSnap{}, err
	}
	mode := compvar.Mode(dto.LiabMode)
	switch mode {
	case compvar.StructMode:
		rec, err := compvar.DataToStructRec(*dto.StructVar)
		if err != nil {
			return ExecSnap{}, err
		}
		return ExecSnap{ref, rec}, nil
	case compvar.LinearMode:
		rec, err := compvar.DataToLinearRec(*dto.LinearVar)
		if err != nil {
			return ExecSnap{}, err
		}
		return ExecSnap{ref, rec}, nil
	default:
		panic(compvar.ErrUnexpectedMode(mode))
	}
}
