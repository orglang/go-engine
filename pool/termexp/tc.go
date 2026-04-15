package termexp

import (
	"github.com/orglang/go-sdk/adt/poolexp"

	"orglang/go-engine/adt/semterm"
	"orglang/go-engine/adt/semtype"
)

func MsgFromExpSpecNilable(spec ExpSpec) *poolexp.ExpSpec {
	if spec == nil {
		return nil
	}
	dto := MsgFromExpSpec(spec)
	return &dto
}

func MsgFromExpSpec(s ExpSpec) poolexp.ExpSpec {
	switch spec := s.(type) {
	case HireSpec:
		return poolexp.ExpSpec{K: poolexp.Hire, Hire: MsgFromHireSpec(spec)}
	case ApplySpec:
		return poolexp.ExpSpec{K: poolexp.Apply, Apply: MsgFromApplySpec(spec)}
	case AcquireSpec:
		return poolexp.ExpSpec{K: poolexp.Acquire, Acquire: MsgFromAcquireSpec(spec)}
	case AcceptSpec:
		return poolexp.ExpSpec{K: poolexp.Accept, Accept: MsgFromAcceptSpec(spec)}
	case ReleaseSpec:
		return poolexp.ExpSpec{K: poolexp.Release, Release: MsgFromReleaseSpec(spec)}
	case DetachSpec:
		return poolexp.ExpSpec{K: poolexp.Detach, Detach: MsgFromDetachSpec(spec)}
	case SpawnSpec2:
		return poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef:  semtype.MsgFromRef(spec.ProcDescRef),
				ProcImplRefs: semterm.MsgFromRefs(spec.ProcImplRefs),
			},
		}
	default:
		panic(ErrSpecTypeUnexpected(s))
	}
}

func MsgToExpSpecNilable(dto *poolexp.ExpSpec) (ExpSpec, error) {
	if dto == nil {
		return nil, nil
	}
	return MsgToExpSpec(*dto)
}

func MsgToExpSpec(dto poolexp.ExpSpec) (ExpSpec, error) {
	switch dto.K {
	case poolexp.Acquire:
		return MsgToAcquireSpec(dto.Acquire)
	case poolexp.Accept:
		return MsgToAcceptSpec(dto.Accept)
	case poolexp.Hire:
		return MsgToHireSpec(dto.Hire)
	case poolexp.Apply:
		return MsgToApplySpec(dto.Apply)
	case poolexp.Release:
		return MsgToReleaseSpec(dto.Release)
	case poolexp.Detach:
		return MsgToDetachSpec(dto.Detach)
	case poolexp.Spawn:
		descRef, err := semtype.MsgToRef(dto.Spawn.ProcDescRef)
		if err != nil {
			return nil, err
		}
		implRefs, err := semterm.MsgToRefs(dto.Spawn.ProcImplRefs)
		if err != nil {
			return nil, err
		}
		return SpawnSpec2{ProcDescRef: descRef, ProcImplRefs: implRefs}, nil
	default:
		panic(poolexp.ErrUnexpectedExpKind(dto.K))
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
