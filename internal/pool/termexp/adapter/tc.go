package adapter

import (
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/termsem"

	termexp "orglang/go-engine/internal/pool/termexp/core"
)

func DataFromExpSpec(s termexp.ExpSpec) termexp.ExpSpecDS {
	switch spec := s.(type) {
	case termexp.AcquireSpec:
		return termexp.ExpSpecDS{K: termexp.AcquireKind, Acquire: DataFromAcquireSpec(spec)}
	case termexp.AcceptSpec:
		return termexp.ExpSpecDS{K: termexp.AcceptKind, Accept: DataFromAcceptSpec(spec)}
	case termexp.HireSpec:
		return termexp.ExpSpecDS{K: termexp.HireKind, Hire: DataFromHireSpec(spec)}
	case termexp.ApplySpec:
		return termexp.ExpSpecDS{K: termexp.ApplyKind, Hire: DataFromApplySpec(spec)}
	case termexp.ReleaseSpec:
		return termexp.ExpSpecDS{K: termexp.ReleaseKind, Release: DataFromReleaseSpec(spec)}
	case termexp.DetachSpec:
		return termexp.ExpSpecDS{K: termexp.DetachKind, Detach: DataFromDetachSpec(spec)}
	default:
		panic(termexp.ErrSpecTypeUnexpected(s))
	}
}

func DataToExpSpec(dto termexp.ExpSpecDS) (termexp.ExpSpec, error) {
	switch dto.K {
	case termexp.AcquireKind:
		return DataToAcquireSpec(dto.Acquire)
	case termexp.AcceptKind:
		return DataToAcceptSpec(dto.Accept)
	case termexp.HireKind:
		return DataToHireSpec(dto.Hire)
	case termexp.ApplyKind:
		return DataToApplySpec(dto.Apply)
	case termexp.ReleaseKind:
		return DataToReleaseSpec(dto.Release)
	case termexp.DetachKind:
		return DataToDetachSpec(dto.Detach)
	default:
		panic(termexp.ErrExpKindUnexpected(dto.K))
	}
}

func DataFromExpRec(r termexp.ExpRec) termexp.ExpRecDS {
	switch rec := r.(type) {
	case termexp.AcquireRec:
		return termexp.ExpRecDS{K: termexp.AcquireKind, Acquire: DataFromAcquireRec(rec)}
	case termexp.AcceptRec:
		return termexp.ExpRecDS{K: termexp.AcceptKind, Accept: DataFromAcceptRec(rec)}
	case termexp.HireRec:
		return termexp.ExpRecDS{K: termexp.HireKind, Hire: DataFromHireRec(rec)}
	case termexp.ApplyRec:
		return termexp.ExpRecDS{K: termexp.ApplyKind, Apply: DataFromApplyRec(rec)}
	case termexp.ReleaseRec:
		return termexp.ExpRecDS{K: termexp.ReleaseKind, Release: DataFromReleaseRec(rec)}
	case termexp.DetachRec:
		return termexp.ExpRecDS{K: termexp.DetachKind, Detach: DataFromDetachRec(rec)}
	default:
		panic(termexp.ErrRecTypeUnexpected(r))
	}
}

func DataToExpRec(dto termexp.ExpRecDS) (termexp.ExpRec, error) {
	switch dto.K {
	case termexp.AcquireKind:
		return DataToAcquireRec(dto.Acquire)
	case termexp.AcceptKind:
		return DataToAcceptRec(dto.Accept)
	case termexp.HireKind:
		return DataToHireRec(dto.Hire)
	case termexp.ApplyKind:
		return DataToApplyRec(dto.Apply)
	case termexp.ReleaseKind:
		return DataToReleaseRec(dto.Release)
	case termexp.DetachKind:
		return DataToDetachRec(dto.Detach)
	default:
		panic(termexp.ErrExpKindUnexpected(dto.K))
	}
}
