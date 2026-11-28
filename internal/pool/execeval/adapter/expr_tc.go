package adapter

import (
	"fmt"

	compexp "orglang/go-engine/internal/pool/execeval/core"
)

func DataFromExprSpec(s compexp.ExprSpec) ExprSpec {
	switch spec := s.(type) {
	case compexp.AcquireSpec:
		return ExprSpec{K: AcquireKind, Acquire: DataFromAcquireSpec(spec)}
	case compexp.AcceptSpec:
		return ExprSpec{K: AcceptKind, Accept: DataFromAcceptSpec(spec)}
	case compexp.HireSpec:
		return ExprSpec{K: HireKind, Hire: DataFromHireSpec(spec)}
	case compexp.ApplySpec:
		return ExprSpec{K: ApplyKind, Apply: DataFromApplySpec(spec)}
	case compexp.ReleaseSpec:
		return ExprSpec{K: ReleaseKind, Release: DataFromReleaseSpec(spec)}
	case compexp.DetachSpec:
		return ExprSpec{K: DetachKind, Detach: DataFromDetachSpec(spec)}
	default:
		panic(compexp.ErrSpecTypeUnexpected(s))
	}
}

func DataToExprSpec(dto ExprSpec) (compexp.ExprSpec, error) {
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
		panic(fmt.Errorf("expr kind unexpected: %v", dto.K))
	}
}

func DataFromExprRec(r compexp.ExprRec) ExprRec {
	switch rec := r.(type) {
	case compexp.AcquireRec:
		return ExprRec{K: AcquireKind, Acquire: DataFromAcquireRec(rec)}
	case compexp.AcceptRec:
		return ExprRec{K: AcceptKind, Accept: DataFromAcceptRec(rec)}
	case compexp.HireRec:
		return ExprRec{K: HireKind, Hire: DataFromHireRec(rec)}
	case compexp.ApplyRec:
		return ExprRec{K: ApplyKind, Apply: DataFromApplyRec(rec)}
	case compexp.ReleaseRec:
		return ExprRec{K: ReleaseKind, Release: DataFromReleaseRec(rec)}
	case compexp.DetachRec:
		return ExprRec{K: DetachKind, Detach: DataFromDetachRec(rec)}
	default:
		panic(compexp.ErrRecTypeUnexpected(r))
	}
}

func DataToExprRec(dto ExprRec) (compexp.ExprRec, error) {
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
		panic(fmt.Errorf("expr kind unexpected: %v", dto.K))
	}
}
