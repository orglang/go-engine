package commturn

import (
	"fmt"

	"orglang/go-engine/proc/termexp"
)

func dataFromStepRec(r TurnRec) (StepRecDS, error) {
	if r == nil {
		return StepRecDS{}, nil
	}
	switch rec := r.(type) {
	case PubRec:
		msgVal, err := termexp.DataFromExpRec(rec.ValExp)
		if err != nil {
			return StepRecDS{}, err
		}
		return StepRecDS{
			K:      msgStep,
			ProcER: msgVal,
		}, nil
	case SubRec:
		svcCont, err := termexp.DataFromExpRec(rec.ContExp)
		if err != nil {
			return StepRecDS{}, err
		}
		return StepRecDS{
			K:      svcStep,
			ProcER: svcCont,
		}, nil
	default:
		panic(ErrRecTypeUnexpected(rec))
	}
}

func dataToStepRec(dto StepRecDS) (TurnRec, error) {
	var nilData StepRecDS
	if dto == nilData {
		return nil, nil
	}
	switch dto.K {
	case msgStep:
		val, err := termexp.DataToExpRec(dto.ProcER)
		if err != nil {
			return nil, err
		}
		return PubRec{ValExp: val}, nil
	case svcStep:
		cont, err := termexp.DataToExpRec(dto.ProcER)
		if err != nil {
			return nil, err
		}
		return SubRec{ContExp: cont}, nil
	default:
		panic(errUnexpectedStepKind(dto.K))
	}
}

func errUnexpectedStepKind(k stepKindDS) error {
	return fmt.Errorf("unexpected step kind: %v", k)
}
