package core

import (
	"orglang/go-engine/internal/adt/comproot"
	compexp "orglang/go-engine/internal/proc/compexec/core"
)

type StepSpec struct {
	CompRef comproot.SemRef
	ProcExp compexp.ExpSpec
}
