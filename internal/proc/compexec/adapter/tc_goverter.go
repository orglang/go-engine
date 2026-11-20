package adapter

import (
	"orglang/go-engine/internal/adt/compsem"
	compsemadapter "orglang/go-engine/internal/adt/compsem/adapter"
	compexec "orglang/go-engine/internal/proc/compexec/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compsem/adapter:Data.*
// goverter:extend orglang/go-engine/internal/adt/compvar/adapter:Data.*
var (
	// goverter:map . CompRef | dataToSemRef
	DataToExecRec func(ExecRecDS) (compexec.ExecRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(compexec.ExecRec) ExecRecDS
	DataFromMod     func(compexec.ExecMod) (ExecModDS, error)
	// goverter:ignore CompRN
	dataToSemRef func(ExecRecDS) (compsem.SemRef, error)
)
