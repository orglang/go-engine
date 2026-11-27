package core

import (
	"orglang/go-engine/internal/adt/comproot"
	connact "orglang/go-engine/internal/proc/connact/core"
)

type StepSpec struct {
	CompRef comproot.SemRef
	ProcExp connact.ExpSpec
}
