package adapter

import (
	sdktermexp "github.com/orglang/go-sdk/pool/termexp"

	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/termsem"

	termexp "orglang/go-engine/pool/termexp/core"
)

func MsgFromExpSpecNilable(spec termexp.ExpSpec) *sdktermexp.ExpSpec {
	if spec == nil {
		return nil
	}
	dto := MsgFromExpSpec(spec)
	return &dto
}

func MsgFromExpSpec(s termexp.ExpSpec) sdktermexp.ExpSpec {
	switch spec := s.(type) {
	case termexp.HireSpec:
		return sdktermexp.ExpSpec{K: sdktermexp.Hire, Hire: MsgFromHireSpec(spec)}
	case termexp.ApplySpec:
		return sdktermexp.ExpSpec{K: sdktermexp.Apply, Apply: MsgFromApplySpec(spec)}
	case termexp.AcquireSpec:
		return sdktermexp.ExpSpec{K: sdktermexp.Acquire, Acquire: MsgFromAcquireSpec(spec)}
	case termexp.AcceptSpec:
		return sdktermexp.ExpSpec{K: sdktermexp.Accept, Accept: MsgFromAcceptSpec(spec)}
	case termexp.ReleaseSpec:
		return sdktermexp.ExpSpec{K: sdktermexp.Release, Release: MsgFromReleaseSpec(spec)}
	case termexp.DetachSpec:
		return sdktermexp.ExpSpec{K: sdktermexp.Detach, Detach: MsgFromDetachSpec(spec)}
	case termexp.SpawnSpec2:
		return sdktermexp.ExpSpec{
			K: sdktermexp.Spawn,
			Spawn: &sdktermexp.SpawnSpec{
				ProcTermRef:  termsem.MsgFromRef(spec.ProcTermRef),
				ProcCompRefs: compsem.MsgFromRefs(spec.ProcCompRefs),
			},
		}
	default:
		panic(termexp.ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpecNilable(dto *sdktermexp.ExpSpec) (termexp.ExpSpec, error) {
	if dto == nil {
		return nil, nil
	}
	return MsgToExpSpec(*dto)
}

func MsgToExpSpec(dto sdktermexp.ExpSpec) (termexp.ExpSpec, error) {
	switch dto.K {
	case sdktermexp.Acquire:
		return MsgToAcquireSpec(dto.Acquire)
	case sdktermexp.Accept:
		return MsgToAcceptSpec(dto.Accept)
	case sdktermexp.Hire:
		return MsgToHireSpec(dto.Hire)
	case sdktermexp.Apply:
		return MsgToApplySpec(dto.Apply)
	case sdktermexp.Release:
		return MsgToReleaseSpec(dto.Release)
	case sdktermexp.Detach:
		return MsgToDetachSpec(dto.Detach)
	case sdktermexp.Spawn:
		termRef, err := termsem.MsgToRef(dto.Spawn.ProcTermRef)
		if err != nil {
			return nil, err
		}
		compRefs, err := compsem.MsgToRefs(dto.Spawn.ProcCompRefs)
		if err != nil {
			return nil, err
		}
		return termexp.SpawnSpec2{ProcTermRef: termRef, ProcCompRefs: compRefs}, nil
	default:
		panic(sdktermexp.ErrUnexpectedExpKind(dto.K))
	}
}

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
