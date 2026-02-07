package procexec

import (
	"github.com/orglang/go-sdk/adt/procexec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/uniqref:Msg.*
var (
	MsgToExecRef   func(procexec.ExecRef) (ExecRef, error)
	MsgFromExecRef func(ExecRef) procexec.ExecRef
	// goverter:ignore ChnlBRs ProcSRs
	MsgToExecSnap   func(procexec.ExecSnap) (ExecSnap, error)
	MsgFromExecSnap func(ExecSnap) procexec.ExecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/uniqref:Data.*
// goverter:extend orglang/go-engine/adt/procbind:Data.*
// goverter:extend orglang/go-engine/adt/procstep:Data.*
var (
	DataFromMod func(ExecMod) (execModDS, error)
)
