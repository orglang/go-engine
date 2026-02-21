package poolstep

import (
	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/poolexp"
	"orglang/go-engine/adt/uniqsym"
)

type StepSpec struct {
	ExecRef descsem.SemRef
	ProcQN  uniqsym.ADT
	ProcES  poolexp.ExpSpec
}
