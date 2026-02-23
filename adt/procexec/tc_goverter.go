package procexec

import (
	"github.com/orglang/go-sdk/adt/procexec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
var (
	// goverter:ignore ChnlVRs ProcSRs
	MsgToExecSnap   func(procexec.ExecSnap) (ExecSnap, error)
	MsgFromExecSnap func(ExecSnap) procexec.ExecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/procstep:Data.*
// goverter:extend orglang/go-engine/adt/implvar:Data.*
var (
	// goverter:autoMap ImplRef
	DataFromExecRec func(ExecRec) execRecDS
	DataFromMod     func(ExecMod) (execModDS, error)
)
