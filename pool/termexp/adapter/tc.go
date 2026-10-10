package adapter

import (
	"github.com/orglang/go-sdk/pool/termexp"

	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/termsem"

	pooltermexp "orglang/go-engine/pool/termexp/core"
)

func MsgFromExpSpecNilable(spec pooltermexp.ExpSpec) *termexp.ExpSpec {
	if spec == nil {
		return nil
	}
	dto := MsgFromExpSpec(spec)
	return &dto
}

func MsgFromExpSpec(s pooltermexp.ExpSpec) termexp.ExpSpec {
	switch spec := s.(type) {
	case pooltermexp.HireSpec:
		return termexp.ExpSpec{K: termexp.Hire, Hire: MsgFromHireSpec(spec)}
	case pooltermexp.ApplySpec:
		return termexp.ExpSpec{K: termexp.Apply, Apply: MsgFromApplySpec(spec)}
	case pooltermexp.AcquireSpec:
		return termexp.ExpSpec{K: termexp.Acquire, Acquire: MsgFromAcquireSpec(spec)}
	case pooltermexp.AcceptSpec:
		return termexp.ExpSpec{K: termexp.Accept, Accept: MsgFromAcceptSpec(spec)}
	case pooltermexp.ReleaseSpec:
		return termexp.ExpSpec{K: termexp.Release, Release: MsgFromReleaseSpec(spec)}
	case pooltermexp.DetachSpec:
		return termexp.ExpSpec{K: termexp.Detach, Detach: MsgFromDetachSpec(spec)}
	case pooltermexp.SpawnSpec2:
		return termexp.ExpSpec{
			K: termexp.Spawn,
			Spawn: &termexp.SpawnSpec{
				ProcTermRef:  termsem.MsgFromRef(spec.ProcTermRef),
				ProcCompRefs: compsem.MsgFromRefs(spec.ProcCompRefs),
			},
		}
	default:
		panic(pooltermexp.ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpecNilable(dto *termexp.ExpSpec) (pooltermexp.ExpSpec, error) {
	if dto == nil {
		return nil, nil
	}
	return MsgToExpSpec(*dto)
}

func MsgToExpSpec(dto termexp.ExpSpec) (pooltermexp.ExpSpec, error) {
	switch dto.K {
	case termexp.Acquire:
		return MsgToAcquireSpec(dto.Acquire)
	case termexp.Accept:
		return MsgToAcceptSpec(dto.Accept)
	case termexp.Hire:
		return MsgToHireSpec(dto.Hire)
	case termexp.Apply:
		return MsgToApplySpec(dto.Apply)
	case termexp.Release:
		return MsgToReleaseSpec(dto.Release)
	case termexp.Detach:
		return MsgToDetachSpec(dto.Detach)
	case termexp.Spawn:
		termRef, err := termsem.MsgToRef(dto.Spawn.ProcTermRef)
		if err != nil {
			return nil, err
		}
		compRefs, err := compsem.MsgToRefs(dto.Spawn.ProcCompRefs)
		if err != nil {
			return nil, err
		}
		return pooltermexp.SpawnSpec2{ProcTermRef: termRef, ProcCompRefs: compRefs}, nil
	default:
		panic(termexp.ErrUnexpectedExpKind(dto.K))
	}
}

func DataFromExpSpec(s pooltermexp.ExpSpec) pooltermexp.ExpSpecDS {
	switch spec := s.(type) {
	case pooltermexp.AcquireSpec:
		return pooltermexp.ExpSpecDS{K: pooltermexp.AcquireKind, Acquire: DataFromAcquireSpec(spec)}
	case pooltermexp.AcceptSpec:
		return pooltermexp.ExpSpecDS{K: pooltermexp.AcceptKind, Accept: DataFromAcceptSpec(spec)}
	case pooltermexp.HireSpec:
		return pooltermexp.ExpSpecDS{K: pooltermexp.HireKind, Hire: DataFromHireSpec(spec)}
	case pooltermexp.ApplySpec:
		return pooltermexp.ExpSpecDS{K: pooltermexp.ApplyKind, Hire: DataFromApplySpec(spec)}
	case pooltermexp.ReleaseSpec:
		return pooltermexp.ExpSpecDS{K: pooltermexp.ReleaseKind, Release: DataFromReleaseSpec(spec)}
	case pooltermexp.DetachSpec:
		return pooltermexp.ExpSpecDS{K: pooltermexp.DetachKind, Detach: DataFromDetachSpec(spec)}
	default:
		panic(pooltermexp.ErrSpecTypeUnexpected(s))
	}
}

func DataToExpSpec(dto pooltermexp.ExpSpecDS) (pooltermexp.ExpSpec, error) {
	switch dto.K {
	case pooltermexp.AcquireKind:
		return DataToAcquireSpec(dto.Acquire)
	case pooltermexp.AcceptKind:
		return DataToAcceptSpec(dto.Accept)
	case pooltermexp.HireKind:
		return DataToHireSpec(dto.Hire)
	case pooltermexp.ApplyKind:
		return DataToApplySpec(dto.Apply)
	case pooltermexp.ReleaseKind:
		return DataToReleaseSpec(dto.Release)
	case pooltermexp.DetachKind:
		return DataToDetachSpec(dto.Detach)
	default:
		panic(pooltermexp.ErrExpKindUnexpected(dto.K))
	}
}

func DataFromExpRec(r pooltermexp.ExpRec) pooltermexp.ExpRecDS {
	switch rec := r.(type) {
	case pooltermexp.AcquireRec:
		return pooltermexp.ExpRecDS{K: pooltermexp.AcquireKind, Acquire: DataFromAcquireRec(rec)}
	case pooltermexp.AcceptRec:
		return pooltermexp.ExpRecDS{K: pooltermexp.AcceptKind, Accept: DataFromAcceptRec(rec)}
	case pooltermexp.HireRec:
		return pooltermexp.ExpRecDS{K: pooltermexp.HireKind, Hire: DataFromHireRec(rec)}
	case pooltermexp.ApplyRec:
		return pooltermexp.ExpRecDS{K: pooltermexp.ApplyKind, Apply: DataFromApplyRec(rec)}
	case pooltermexp.ReleaseRec:
		return pooltermexp.ExpRecDS{K: pooltermexp.ReleaseKind, Release: DataFromReleaseRec(rec)}
	case pooltermexp.DetachRec:
		return pooltermexp.ExpRecDS{K: pooltermexp.DetachKind, Detach: DataFromDetachRec(rec)}
	default:
		panic(pooltermexp.ErrRecTypeUnexpected(r))
	}
}

func DataToExpRec(dto pooltermexp.ExpRecDS) (pooltermexp.ExpRec, error) {
	switch dto.K {
	case pooltermexp.AcquireKind:
		return DataToAcquireRec(dto.Acquire)
	case pooltermexp.AcceptKind:
		return DataToAcceptRec(dto.Accept)
	case pooltermexp.HireKind:
		return DataToHireRec(dto.Hire)
	case pooltermexp.ApplyKind:
		return DataToApplyRec(dto.Apply)
	case pooltermexp.ReleaseKind:
		return DataToReleaseRec(dto.Release)
	case pooltermexp.DetachKind:
		return DataToDetachRec(dto.Detach)
	default:
		panic(pooltermexp.ErrExpKindUnexpected(dto.K))
	}
}
