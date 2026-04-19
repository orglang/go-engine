package compstep

import (
	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/proc/termexp"
)

type StepSpec struct {
	CompRef compsem.SemRef
	ProcExp termexp.ExpSpec
}
