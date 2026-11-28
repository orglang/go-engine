package core

import (
	"orglang/go-engine/internal/adt/execroot"
	commact "orglang/go-engine/internal/proc/commact/core"
)

type StepSpec struct {
	RootRef execroot.RootRef
	ProcExp commact.ExpSpec
}
