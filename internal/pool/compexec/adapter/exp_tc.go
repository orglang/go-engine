package adapter

import (
	"fmt"

	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"

	compexp "orglang/go-engine/internal/pool/compexec/exp"
)

func DataFromCompSpec(s compexp.ExpSpec) CompSpec {
	switch spec := s.(type) {
	case compexp.AcquireSpec:
		return CompSpec{K: AcquireKind, Acquire: DataFromAcquireSpec(spec)}
	case compexp.AcceptSpec:
		return CompSpec{K: AcceptKind, Accept: DataFromAcceptSpec(spec)}
	case compexp.HireSpec:
		return CompSpec{K: HireKind, Hire: DataFromHireSpec(spec)}
	case compexp.ApplySpec:
		return CompSpec{K: ApplyKind, Apply: DataFromApplySpec(spec)}
	case compexp.ReleaseSpec:
		return CompSpec{K: ReleaseKind, Release: DataFromReleaseSpec(spec)}
	case compexp.DetachSpec:
		return CompSpec{K: DetachKind, Detach: DataFromDetachSpec(spec)}
	default:
		panic(compexp.ErrSpecTypeUnexpected(s))
	}
}

func DataToCompSpec(dto CompSpec) (compexp.ExpSpec, error) {
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

func DataFromExpRec(r compexp.ExpRec) ExpRec {
	switch rec := r.(type) {
	case compexp.AcquireRec:
		return ExpRec{K: AcquireKind, Acquire: DataFromAcquireRec(rec)}
	case compexp.AcceptRec:
		return ExpRec{K: AcceptKind, Accept: DataFromAcceptRec(rec)}
	case compexp.HireRec:
		return ExpRec{K: HireKind, Hire: DataFromHireRec(rec)}
	case compexp.ApplyRec:
		return ExpRec{K: ApplyKind, Apply: DataFromApplyRec(rec)}
	case compexp.ReleaseRec:
		return ExpRec{K: ReleaseKind, Release: DataFromReleaseRec(rec)}
	case compexp.DetachRec:
		return ExpRec{K: DetachKind, Detach: DataFromDetachRec(rec)}
	default:
		panic(compexp.ErrRecTypeUnexpected(r))
	}
}

func DataToExpRec(dto ExpRec) (compexp.ExpRec, error) {
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
		panic(fmt.Errorf("exp kind unexpected: %v", dto.K))
	}
}

