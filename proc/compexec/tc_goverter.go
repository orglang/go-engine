package compexec

import (
	"orglang/go-engine/adt/compsem"

	"github.com/orglang/go-sdk/proc/compexec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
var (
	// goverter:ignore LinearVars
	MsgToExecSnap   func(compexec.ExecSnap) (ExecSnap, error)
	MsgFromExecSnap func(ExecSnap) compexec.ExecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/proc/commturn:Data.*
// goverter:extend orglang/go-engine/adt/compvar:Data.*
var (
	// goverter:map . CompRef | dataToSemRef
	DataToExecRec func(execRecDS) (ExecRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(ExecRec) execRecDS
	DataFromMod     func(ExecMod) (execModDS, error)
	// goverter:ignore CompRN
	dataToSemRef func(execRecDS) (compsem.SemRef, error)
)
