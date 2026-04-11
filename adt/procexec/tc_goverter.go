package procexec

import (
	"orglang/go-engine/adt/implsem"

	"github.com/orglang/go-sdk/adt/procexec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
var (
	// goverter:ignore LinearVars
	MsgToExecSnap   func(procexec.ExecSnap) (ExecSnap, error)
	MsgFromExecSnap func(ExecSnap) procexec.ExecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/procstep:Data.*
// goverter:extend orglang/go-engine/adt/implvar:Data.*
var (
	// goverter:map . ImplRef | DataToImplRef
	DataToExecRec func(execRecDS) (ExecRec, error)
	// goverter:autoMap ImplRef
	DataFromExecRec func(ExecRec) execRecDS
	DataFromMod     func(ExecMod) (execModDS, error)
	// goverter:ignore ImplRN
	DataToImplRef func(execRecDS) (implsem.SemRef, error)
)
