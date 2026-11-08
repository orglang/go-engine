package adapter

import (
	"orglang/go-engine/internal/adt/compsem"
	compexec "orglang/go-engine/internal/proc/compexec/core"

	sdkcompexec "github.com/orglang/go-sdk/proc/compexec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
var (
	// goverter:ignore LinearVars
	MsgToExecSnap   func(sdkcompexec.ExecSnap) (compexec.ExecSnap, error)
	MsgFromExecSnap func(compexec.ExecSnap) sdkcompexec.ExecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/proc/commturn/adapter:Data.*
// goverter:extend orglang/go-engine/internal/adt/compvar:Data.*
var (
	// goverter:map . CompRef | dataToSemRef
	DataToExecRec func(compexec.ExecRecDS) (compexec.ExecRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(compexec.ExecRec) compexec.ExecRecDS
	DataFromMod     func(compexec.ExecMod) (compexec.ExecModDS, error)
	// goverter:ignore CompRN
	dataToSemRef func(compexec.ExecRecDS) (compsem.SemRef, error)
)
