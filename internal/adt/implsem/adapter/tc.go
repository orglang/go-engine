package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/implsem"
	"orglang/go-engine/internal/adt/uniqsym"
)

func dataFromRec(rec implsem.SemRec) semRecDS {
	return semRecDS{
		ImplQN: uniqsym.ConvertToString(rec.ImplQN),
		ImplID: identity.ConvertToString(rec.ImplID),
		Kind:   int16(rec.Kind),
	}
}
