package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/implsem"
	"orglang/go-engine/internal/adt/uniqsym"
)

func dataFromRec(rec implsem.SemRec) (semRecDS, error) {
	return semRecDS{
		ImplQN: uniqsym.ConvertToString(rec.ImplQN),
		ImplID: identity.ConvertToString(rec.ImplID),
		Kind:   int16(rec.Kind),
	}, nil
}

func dataToRec(dto semRecDS) (implsem.SemRec, error) {
	qn, err := uniqsym.ConvertFromString(dto.ImplQN)
	if err != nil {
		return implsem.SemRec{}, err
	}
	id, err := identity.ConvertFromString(dto.ImplID)
	if err != nil {
		return implsem.SemRec{}, err
	}
	return implsem.SemRec{ImplQN: qn, ImplID: id, Kind: implsem.ImplKindFromInt16(dto.Kind)}, nil
}
