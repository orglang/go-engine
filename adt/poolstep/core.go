package poolstep

import (
	"orglang/go-engine/adt/poolexp"
	"orglang/go-engine/adt/uniqref"
	"orglang/go-engine/adt/uniqsym"
)

type StepSpec struct {
	ExecRef uniqref.ADT
	ProcQN  uniqsym.ADT
	ProcES  poolexp.ExpSpec
}
