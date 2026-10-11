package adapter

import (
	"fmt"

	proccommturn "orglang/go-engine/proc/commturn/core"
	termexp "orglang/go-engine/proc/termexp/adapter"
)

func dataFromStepRec(r proccommturn.TurnRec) (proccommturn.StepRecDS, error) {
	if r == nil {
		return proccommturn.StepRecDS{}, nil
	}
	switch rec := r.(type) {
	case proccommturn.PubRec:
		msgVal, err := termexp.DataFromExpRec(rec.ValExp)
		if err != nil {
			return proccommturn.StepRecDS{}, err
		}
		return proccommturn.StepRecDS{
			K:      proccommturn.MsgStep,
			ProcER: msgVal,
		}, nil
	case proccommturn.SubRec:
		svcCont, err := termexp.DataFromExpRec(rec.ContExp)
		if err != nil {
			return proccommturn.StepRecDS{}, err
		}
		return proccommturn.StepRecDS{
			K:      proccommturn.SvcStep,
			ProcER: svcCont,
		}, nil
	default:
		panic(proccommturn.ErrRecTypeUnexpected(rec))
	}
}

func dataToStepRec(dto proccommturn.StepRecDS) (proccommturn.TurnRec, error) {
	var nilData proccommturn.StepRecDS
	if dto == nilData {
		return nil, nil
	}
	switch dto.K {
	case proccommturn.MsgStep:
		val, err := termexp.DataToExpRec(dto.ProcER)
		if err != nil {
			return nil, err
		}
		return proccommturn.PubRec{ValExp: val}, nil
	case proccommturn.SvcStep:
		cont, err := termexp.DataToExpRec(dto.ProcER)
		if err != nil {
			return nil, err
		}
		return proccommturn.SubRec{ContExp: cont}, nil
	default:
		panic(errUnexpectedStepKind(dto.K))
	}
}

func errUnexpectedStepKind(k proccommturn.StepKindDS) error {
	return fmt.Errorf("unexpected step kind: %v", k)
}
