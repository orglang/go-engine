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
	case SpawnSpec2:
		return poolexp.ExpSpec{
			K: poolexp.Spawn,
			Spawn: &poolexp.SpawnSpec{
				ProcDescRef:  descsem.MsgFromRef(spec.ProcDescRef),
				ProcImplRefs: implsem.MsgFromRefs(spec.ProcImplRefs),
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
		commPH, err := symbol.ConvertFromString(dto.Acquire.CommChnlPH)
		if err != nil {
			return nil, err
		}
		contExp, err := MsgToExpSpec(dto.Acquire.ContExp)
		if err != nil {
			return nil, err
		}
		return AcquireSpec{CommChnlPH: commPH, ContExp: contExp}, nil
	case poolexp.Accept:
		commPH, err := symbol.ConvertFromString(dto.Accept.CommChnlPH)
		if err != nil {
			return nil, err
		}
		contExp, err := MsgToExpSpec(dto.Accept.ContExp)
		if err != nil {
			return nil, err
		}
		return AcceptSpec{CommChnlPH: commPH, ContExp: contExp}, nil
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
	default:
		panic(ErrSpecTypeUnexpected(s))
	}
}

func DataToExpSpec(dto ExpSpecDS) (ExpSpec, error) {
	switch dto.K {
	case acquireKind:
		spec, err := DataToAcquireSpec(dto.Acquire)
		if err != nil {
			return nil, err
		}
		return spec, nil
	case acceptKind:
		spec, err := DataToAcceptSpec(dto.Accept)
		if err != nil {
			return nil, err
		}
		return spec, nil
	case hireKind:
		spec, err := DataToHireSpec(dto.Hire)
		if err != nil {
			return nil, err
		}
		return spec, nil
	case applyKind:
		spec, err := DataToApplySpec(dto.Apply)
		if err != nil {
			return nil, err
		}
		return spec, nil
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
	default:
		panic(ErrRecTypeUnexpected(r))
	}
}

func DataToExpRec(dto ExpRecDS) (ExpRec, error) {
	switch dto.K {
	case acquireKind:
		rec, err := DataToAcquireRec(dto.Acquire)
		if err != nil {
			return nil, err
		}
		return rec, nil
	case acceptKind:
		rec, err := DataToAcceptRec(dto.Accept)
		if err != nil {
			return nil, err
		}
		return rec, nil
	default:
		panic(ErrExpKindUnexpected(dto.K))
	}
}
