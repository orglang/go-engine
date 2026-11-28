package core

import (
	"orglang/go-engine/internal/adt/execroot"
	connact "orglang/go-engine/internal/proc/connact/core"
)

type StepSpec struct {
	RootRef execroot.RootRef
	ProcExp connact.ExpSpec
}
