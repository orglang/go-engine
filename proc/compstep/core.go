package compstep

import (
	"orglang/go-engine/adt/semterm"

	"orglang/go-engine/proc/termexp"
)

type StepSpec struct {
	CompRef semterm.TermRef
	ProcExp termexp.ExpSpec
}
