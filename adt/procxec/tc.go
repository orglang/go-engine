package procxec

import (
	"fmt"

	"orglang/orglang/adt/procdef"
)

func dataFromSemRec(r SemRec) (SemRecDS, error) {
	if r == nil {
		return SemRecDS{}, nil
	}
	switch rec := r.(type) {
	case MsgRec:
		msgVal, err := procdef.DataFromTermRec(rec.Val)
		if err != nil {
			return SemRecDS{}, err
		}
		return SemRecDS{
			K:  msgKind,
			TR: msgVal,
		}, nil
	case SvcRec:
		svcCont, err := procdef.DataFromTermRec(rec.Cont)
		if err != nil {
			return SemRecDS{}, err
		}
		return SemRecDS{
			K:  svcKind,
			TR: svcCont,
		}, nil
	default:
		panic(ErrRootTypeUnexpected(rec))
	}
}

func dataToSemRec(dto SemRecDS) (SemRec, error) {
	var nilData SemRecDS
	if dto == nilData {
		return nil, nil
	}
	switch dto.K {
	case msgKind:
		val, err := procdef.DataToTermRec(dto.TR)
		if err != nil {
			return nil, err
		}
		return MsgRec{Val: val}, nil
	case svcKind:
		cont, err := procdef.DataToTermRec(dto.TR)
		if err != nil {
			return nil, err
		}
		return SvcRec{Cont: cont}, nil
	default:
		panic(errUnexpectedStepKind(dto.K))
	}
}

func errUnexpectedStepKind(k semKind) error {
	return fmt.Errorf("unexpected step kind: %v", k)
}
