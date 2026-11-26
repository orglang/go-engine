package core

import "orglang/go-engine/internal/adt/comproot"

// ExpSpec is the expression contract required to schedule a computation step.
type ExpSpec interface {
	IsPoolExp()
}

type StepSpec struct {
	CompRef comproot.SemRef
	PoolExp ExpSpec
}
