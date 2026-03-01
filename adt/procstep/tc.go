package procstep

import (
	"fmt"

	"orglang/go-engine/adt/procexp"
)

func dataFromStepRec(r CommRec) (StepRecDS, error) {
	if r == nil {
		return StepRecDS{}, nil
	}
	switch rec := r.(type) {
	case PubRec:
		msgVal, err := procexp.DataFromExpRec(rec.ValExp)
		if err != nil {
			return StepRecDS{}, err
		}
		return StepRecDS{
			K:      msgStep,
			ProcER: msgVal,
		}, nil
	case SubRec:
		svcCont, err := procexp.DataFromExpRec(rec.ContExp)
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

func dataToStepRec(dto StepRecDS) (CommRec, error) {
	var nilData StepRecDS
	if dto == nilData {
		return nil, nil
	}
	switch dto.K {
	case msgStep:
		val, err := procexp.DataToExpRec(dto.ProcER)
		if err != nil {
			return nil, err
		}
		return PubRec{ValExp: val}, nil
	case svcStep:
		cont, err := procexp.DataToExpRec(dto.ProcER)
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
