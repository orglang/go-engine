package core

import (
	"orglang/go-engine/internal/adt/compsem"
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

type StepSpec struct {
	CompRef compsem.SemRef
	ProcExp termexp.ExpSpec
}
