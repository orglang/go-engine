package core

import "orglang/go-engine/internal/adt/execroot"

// ExprSpec is the expression contract required to schedule an execution step.
type ExprSpec interface {
	IsPoolExpr()
}

type StepSpec struct {
	RootRef  execroot.RootRef
	PoolExpr ExprSpec
}
