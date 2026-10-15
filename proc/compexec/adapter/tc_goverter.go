package adapter

import (
	"orglang/go-engine/adt/compsem"
	core "orglang/go-engine/proc/compexec/core"

	"github.com/orglang/go-sdk/proc/compexec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
var (
	// goverter:ignore LinearVars
	MsgToExecSnap   func(compexec.ExecSnap) (core.ExecSnap, error)
	MsgFromExecSnap func(core.ExecSnap) compexec.ExecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/proc/commturn/adapter:Data.*
// goverter:extend orglang/go-engine/adt/compvar:Data.*
var (
	// goverter:map . CompRef | dataToSemRef
	DataToExecRec func(core.ExecRecDS) (core.ExecRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(core.ExecRec) core.ExecRecDS
	DataFromMod     func(core.ExecMod) (core.ExecModDS, error)
	// goverter:ignore CompRN
	dataToSemRef func(core.ExecRecDS) (compsem.SemRef, error)
)
