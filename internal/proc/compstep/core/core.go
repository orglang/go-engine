package core

import (
	"orglang/go-engine/internal/adt/comproot"
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

type StepSpec struct {
	CompRef comproot.SemRef
	ProcExp termexp.CompSpec
}
