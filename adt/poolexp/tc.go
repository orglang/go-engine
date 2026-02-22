package poolexp

import (
	"github.com/orglang/go-sdk/adt/poolexp"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
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
		return poolexp.ExpSpec{
			K: poolexp.Hire,
			Hire: &poolexp.HireSpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				ProcDescQN: uniqsym.ConvertToString(spec.ProcDescQN),
			},
		}
	case ApplySpec:
		return poolexp.ExpSpec{
			K: poolexp.Apply,
			Apply: &poolexp.ApplySpec{
				CommChnlPH: symbol.ConvertToString(spec.CommChnlPH),
				ProcDescQN: uniqsym.ConvertToString(spec.ProcDescQN),
			},
		}
	case SpawnSpec:
		return poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef:  descsem.MsgFromRef(spec.ProcDescRef),
				ProcImplRefs: implsem.MsgFromRefs(spec.ProcImplRefs),
			},
		}
	default:
		panic(ErrExpTypeUnexpected(s))
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
	case poolexp.Hire:
		commPH, err := symbol.ConvertFromString(dto.Hire.CommChnlPH)
		if err != nil {
			return nil, err
		}
		descQN, err := uniqsym.ConvertFromString(dto.Hire.ProcDescQN)
		if err != nil {
			return nil, err
		}
		return HireSpec{CommChnlPH: commPH, ProcDescQN: descQN}, nil
	case poolexp.Apply:
		commPH, err := symbol.ConvertFromString(dto.Apply.CommChnlPH)
		if err != nil {
			return nil, err
		}
		descQN, err := uniqsym.ConvertFromString(dto.Apply.ProcDescQN)
		if err != nil {
			return nil, err
		}
		return ApplySpec{CommChnlPH: commPH, ProcDescQN: descQN}, nil
	case poolexp.Spawn:
		descRef, err := descsem.MsgToRef(dto.Spawn.ProcDescRef)
		if err != nil {
			return nil, err
		}
		implRefs, err := implsem.MsgToRefs(dto.Spawn.ProcImplRefs)
		if err != nil {
			return nil, err
		}
		return SpawnSpec{ProcDescRef: descRef, ProcImplRefs: implRefs}, nil
	default:
		panic(poolexp.ErrUnexpectedExpKind(dto.K))
	}
}
