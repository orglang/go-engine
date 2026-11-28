package adapter

import (
	"fmt"

	compexpr "orglang/go-engine/internal/pool/execeval/core"
)

func DataFromExprSpec(s compexpr.ExprSpec) ExprSpec {
	switch spec := s.(type) {
	case compexpr.AcquireSpec:
		return ExprSpec{K: AcquireKind, Acquire: DataFromAcquireSpec(spec)}
	case compexpr.AcceptSpec:
		return ExprSpec{K: AcceptKind, Accept: DataFromAcceptSpec(spec)}
	case compexpr.HireSpec:
		return ExprSpec{K: HireKind, Hire: DataFromHireSpec(spec)}
	case compexpr.ApplySpec:
		return ExprSpec{K: ApplyKind, Apply: DataFromApplySpec(spec)}
	case compexpr.ReleaseSpec:
		return ExprSpec{K: ReleaseKind, Release: DataFromReleaseSpec(spec)}
	case compexpr.DetachSpec:
		return ExprSpec{K: DetachKind, Detach: DataFromDetachSpec(spec)}
	default:
		panic(compexpr.ErrSpecTypeUnexpected(s))
	}
}

func DataToExprSpec(dto ExprSpec) (compexpr.ExprSpec, error) {
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

func DataFromExprRec(r compexpr.ExprRec) ExprRec {
	switch rec := r.(type) {
	case compexpr.AcquireRec:
		return ExprRec{K: AcquireKind, Acquire: DataFromAcquireRec(rec)}
	case compexpr.AcceptRec:
		return ExprRec{K: AcceptKind, Accept: DataFromAcceptRec(rec)}
	case compexpr.HireRec:
		return ExprRec{K: HireKind, Hire: DataFromHireRec(rec)}
	case compexpr.ApplyRec:
		return ExprRec{K: ApplyKind, Apply: DataFromApplyRec(rec)}
	case compexpr.ReleaseRec:
		return ExprRec{K: ReleaseKind, Release: DataFromReleaseRec(rec)}
	case compexpr.DetachRec:
		return ExprRec{K: DetachKind, Detach: DataFromDetachRec(rec)}
	default:
		panic(compexpr.ErrRecTypeUnexpected(r))
	}
}

func DataToExprRec(dto ExprRec) (compexpr.ExprRec, error) {
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
