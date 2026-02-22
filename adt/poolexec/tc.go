package poolexec

import (
	"orglang/go-engine/adt/implsem"
)

func ConvertRecToRef(rec ExecRec) implsem.SemRef {
	return rec.ImplRef
}
