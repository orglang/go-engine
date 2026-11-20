package adapter

import (
	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
)

func dataFromRec(rec descsem.SemRec) (semRecDS, error) {
	return semRecDS{
		DescQN: uniqsym.ConvertToString(rec.DescQN),
		DescID: identity.ConvertToString(rec.DescID),
		Kind:   int16(rec.Kind),
	}, nil
}

func dataToRec(dto semRecDS) (descsem.SemRec, error) {
	qn, err := uniqsym.ConvertFromString(dto.DescQN)
	if err != nil {
		return descsem.SemRec{}, err
	}
	id, err := identity.ConvertFromString(dto.DescID)
	if err != nil {
		return descsem.SemRec{}, err
	}
	return descsem.SemRec{DescQN: qn, DescID: id, Kind: descsem.DescKindFromInt16(dto.Kind)}, nil
}
