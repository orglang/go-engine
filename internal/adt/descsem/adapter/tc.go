package adapter

import (
	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
)

func dataFromRec(rec descsem.SemRec) semRecDS {
	return semRecDS{
		DescQN: uniqsym.ConvertToString(rec.DescQN),
		DescID: identity.ConvertToString(rec.DescID),
		Kind:   int16(rec.Kind),
	}
}
