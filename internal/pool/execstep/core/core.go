package core

import "orglang/go-engine/internal/adt/execroot"

// ExpSpec is the expression contract required to schedule an execution step.
type ExpSpec interface {
	IsPoolExp()
}

type StepSpec struct {
	RootRef execroot.RootRef
	PoolExpr ExpSpec
}
