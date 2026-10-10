package compstep

import (
	"orglang/go-engine/adt/compsem"
	termexp "orglang/go-engine/proc/termexp/core"
)

type StepSpec struct {
	CompRef compsem.SemRef
	ProcExp termexp.ExpSpec
}
