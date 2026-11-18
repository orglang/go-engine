package adapter

import (
	"orglang/go-engine/internal/adt/compsem"
	compexec "orglang/go-engine/internal/proc/compexec/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
var (
	// goverter:map . CompRef | dataToSemRef
	DataToExecRec func(compexec.ExecRecDS) (compexec.ExecRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(compexec.ExecRec) compexec.ExecRecDS
	DataFromMod     func(compexec.ExecMod) (compexec.ExecModDS, error)
	// goverter:ignore CompRN
	dataToSemRef func(compexec.ExecRecDS) (compsem.SemRef, error)
)
