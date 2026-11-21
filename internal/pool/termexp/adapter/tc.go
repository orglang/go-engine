package adapter

import (
	"fmt"

	termexp "orglang/go-engine/internal/pool/termexp/core"
)

func DataFromExpSpec(s termexp.ExpSpec) ExpSpec {
	switch spec := s.(type) {
	case termexp.AcquireSpec:
		return ExpSpec{K: AcquireKind, Acquire: DataFromAcquireSpec(spec)}
	case termexp.AcceptSpec:
		return ExpSpec{K: AcceptKind, Accept: DataFromAcceptSpec(spec)}
	case termexp.HireSpec:
		return ExpSpec{K: HireKind, Hire: DataFromHireSpec(spec)}
	case termexp.ApplySpec:
		return ExpSpec{K: ApplyKind, Hire: DataFromApplySpec(spec)}
	case termexp.ReleaseSpec:
		return ExpSpec{K: ReleaseKind, Release: DataFromReleaseSpec(spec)}
	case termexp.DetachSpec:
		return ExpSpec{K: DetachKind, Detach: DataFromDetachSpec(spec)}
	default:
		panic(termexp.ErrSpecTypeUnexpected(s))
	}
}

func DataToExpSpec(dto ExpSpec) (termexp.ExpSpec, error) {
	switch dto.K {
	case AcquireKind:
		return DataToAcquireSpec(dto.Acquire)
	case AcceptKind:
		return DataToAcceptSpec(dto.Accept)
	case HireKind:
		return DataToHireSpec(dto.Hire)
	case ApplyKind:
		return DataToApplySpec(dto.Apply)
	case ReleaseKind:
		return DataToReleaseSpec(dto.Release)
	case DetachKind:
		return DataToDetachSpec(dto.Detach)
	default:
		panic(fmt.Errorf("exp kind unexpected: %v", dto.K))
	}
}

func DataFromExpRec(r termexp.ExpRec) ExpRec {
	switch rec := r.(type) {
	case termexp.AcquireRec:
		return ExpRec{K: AcquireKind, Acquire: DataFromAcquireRec(rec)}
	case termexp.AcceptRec:
		return ExpRec{K: AcceptKind, Accept: DataFromAcceptRec(rec)}
	case termexp.HireRec:
		return ExpRec{K: HireKind, Hire: DataFromHireRec(rec)}
	case termexp.ApplyRec:
		return ExpRec{K: ApplyKind, Apply: DataFromApplyRec(rec)}
	case termexp.ReleaseRec:
		return ExpRec{K: ReleaseKind, Release: DataFromReleaseRec(rec)}
	case termexp.DetachRec:
		return ExpRec{K: DetachKind, Detach: DataFromDetachRec(rec)}
	default:
		panic(termexp.ErrRecTypeUnexpected(r))
	}
}

func DataToExpRec(dto ExpRec) (termexp.ExpRec, error) {
	switch dto.K {
	case AcquireKind:
		return DataToAcquireRec(dto.Acquire)
	case AcceptKind:
		return DataToAcceptRec(dto.Accept)
	case HireKind:
		return DataToHireRec(dto.Hire)
	case ApplyKind:
		return DataToApplyRec(dto.Apply)
	case ReleaseKind:
		return DataToReleaseRec(dto.Release)
	case DetachKind:
		return DataToDetachRec(dto.Detach)
	default:
		panic(termexp.ErrExpKindUnexpected(dto.K))
	}
}
