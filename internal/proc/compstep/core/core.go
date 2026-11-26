package core

import (
	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/symbol"
)

// ExpSpec is the expression contract required to schedule a process step.
type ExpSpec interface {
	Via() symbol.ADT
}

type StepSpec struct {
	CompRef comproot.SemRef
	ProcExp ExpSpec
}
