package adapter

import (
	"fmt"

	commturn "orglang/go-engine/proc/commturn/core"
	termexp "orglang/go-engine/proc/termexp/adapter"
)

func dataFromStepRec(r commturn.TurnRec) (commturn.StepRecDS, error) {
	if r == nil {
		return commturn.StepRecDS{}, nil
	}
	switch rec := r.(type) {
	case commturn.PubRec:
		msgVal, err := termexp.DataFromExpRec(rec.ValExp)
		if err != nil {
			return commturn.StepRecDS{}, err
		}
		return commturn.StepRecDS{
			K:      commturn.MsgStep,
			ProcER: msgVal,
		}, nil
	case commturn.SubRec:
		svcCont, err := termexp.DataFromExpRec(rec.ContExp)
		if err != nil {
			return commturn.StepRecDS{}, err
		}
		return commturn.StepRecDS{
			K:      commturn.SvcStep,
			ProcER: svcCont,
		}, nil
	default:
		panic(commturn.ErrRecTypeUnexpected(rec))
	}
}

func dataToStepRec(dto commturn.StepRecDS) (commturn.TurnRec, error) {
	var nilData commturn.StepRecDS
	if dto == nilData {
		return nil, nil
	}
	switch dto.K {
	case commturn.MsgStep:
		val, err := termexp.DataToExpRec(dto.ProcER)
		if err != nil {
			return nil, err
		}
		return commturn.PubRec{ValExp: val}, nil
	case commturn.SvcStep:
		cont, err := termexp.DataToExpRec(dto.ProcER)
		if err != nil {
			return nil, err
		}
		return commturn.SubRec{ContExp: cont}, nil
	default:
		panic(errUnexpectedStepKind(dto.K))
	}
}

func errUnexpectedStepKind(k commturn.StepKindDS) error {
	return fmt.Errorf("unexpected step kind: %v", k)
}
