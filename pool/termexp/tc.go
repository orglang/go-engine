package termexp

import (
	"github.com/orglang/go-sdk/pool/termexp"

	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/termsem"
)

func MsgFromExpSpecNilable(spec ExpSpec) *termexp.ExpSpec {
	if spec == nil {
		return nil
	}
	dto := MsgFromExpSpec(spec)
	return &dto
}

func MsgFromExpSpec(s ExpSpec) termexp.ExpSpec {
	switch spec := s.(type) {
	case HireSpec:
		return termexp.ExpSpec{K: termexp.Hire, Hire: MsgFromHireSpec(spec)}
	case ApplySpec:
		return termexp.ExpSpec{K: termexp.Apply, Apply: MsgFromApplySpec(spec)}
	case AcquireSpec:
		return termexp.ExpSpec{K: termexp.Acquire, Acquire: MsgFromAcquireSpec(spec)}
	case AcceptSpec:
		return termexp.ExpSpec{K: termexp.Accept, Accept: MsgFromAcceptSpec(spec)}
	case ReleaseSpec:
		return termexp.ExpSpec{K: termexp.Release, Release: MsgFromReleaseSpec(spec)}
	case DetachSpec:
		return termexp.ExpSpec{K: termexp.Detach, Detach: MsgFromDetachSpec(spec)}
	case SpawnSpec2:
		return termexp.ExpSpec{
			K: termexp.Spawn,
			Spawn: &termexp.SpawnSpec{
				ProcTermRef:  termsem.MsgFromRef(spec.ProcTermRef),
				ProcCompRefs: compsem.MsgFromRefs(spec.ProcCompRefs),
			},
		}
	default:
		panic(ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpecNilable(dto *termexp.ExpSpec) (ExpSpec, error) {
	if dto == nil {
		return nil, nil
	}
	return MsgToExpSpec(*dto)
}

func MsgToExpSpec(dto termexp.ExpSpec) (ExpSpec, error) {
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
		return SpawnSpec2{ProcTermRef: termRef, ProcCompRefs: compRefs}, nil
	default:
		panic(termexp.ErrUnexpectedExpKind(dto.K))
	}
}

func DataFromExpSpec(s ExpSpec) ExpSpecDS {
	switch spec := s.(type) {
	case AcquireSpec:
		return ExpSpecDS{K: acquireKind, Acquire: DataFromAcquireSpec(spec)}
	case AcceptSpec:
		return ExpSpecDS{K: acceptKind, Accept: DataFromAcceptSpec(spec)}
	case HireSpec:
		return ExpSpecDS{K: hireKind, Hire: DataFromHireSpec(spec)}
	case ApplySpec:
		return ExpSpecDS{K: applyKind, Hire: DataFromApplySpec(spec)}
	case ReleaseSpec:
		return ExpSpecDS{K: releaseKind, Release: DataFromReleaseSpec(spec)}
	case DetachSpec:
		return ExpSpecDS{K: detachKind, Detach: DataFromDetachSpec(spec)}
	default:
		panic(ErrSpecTypeUnexpected(s))
	}
}

func DataToExpSpec(dto ExpSpecDS) (ExpSpec, error) {
	switch dto.K {
	case acquireKind:
		return DataToAcquireSpec(dto.Acquire)
	case acceptKind:
		return DataToAcceptSpec(dto.Accept)
	case hireKind:
		return DataToHireSpec(dto.Hire)
	case applyKind:
		return DataToApplySpec(dto.Apply)
	case releaseKind:
		return DataToReleaseSpec(dto.Release)
	case detachKind:
		return DataToDetachSpec(dto.Detach)
	default:
		panic(ErrExpKindUnexpected(dto.K))
	}
}

func DataFromExpRec(r ExpRec) ExpRecDS {
	switch rec := r.(type) {
	case AcquireRec:
		return ExpRecDS{K: acquireKind, Acquire: DataFromAcquireRec(rec)}
	case AcceptRec:
		return ExpRecDS{K: acceptKind, Accept: DataFromAcceptRec(rec)}
	case HireRec:
		return ExpRecDS{K: hireKind, Hire: DataFromHireRec(rec)}
	case ApplyRec:
		return ExpRecDS{K: applyKind, Apply: DataFromApplyRec(rec)}
	case ReleaseRec:
		return ExpRecDS{K: releaseKind, Release: DataFromReleaseRec(rec)}
	case DetachRec:
		return ExpRecDS{K: detachKind, Detach: DataFromDetachRec(rec)}
	default:
		panic(ErrRecTypeUnexpected(r))
	}
}

func DataToExpRec(dto ExpRecDS) (ExpRec, error) {
	switch dto.K {
	case acquireKind:
		return DataToAcquireRec(dto.Acquire)
	case acceptKind:
		return DataToAcceptRec(dto.Accept)
	case hireKind:
		return DataToHireRec(dto.Hire)
	case applyKind:
		return DataToApplyRec(dto.Apply)
	case releaseKind:
		return DataToReleaseRec(dto.Release)
	case detachKind:
		return DataToDetachRec(dto.Detach)
	default:
		panic(ErrExpKindUnexpected(dto.K))
	}
}
