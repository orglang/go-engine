package adapter

import (
	"orglang/go-engine/adt/uniqsym"

	compexec "orglang/go-engine/pool/compexec/core"

	sdkcompexec "github.com/orglang/go-sdk/pool/compexec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/compsem:Msg.*
var (
	MsgToExecSpec   func(sdkcompexec.ExecSpec) (compexec.ExecSpec, error)
	MsgFromExecSpec func(compexec.ExecSpec) sdkcompexec.ExecSpec
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/compvar:Convert.*
// goverter:extend orglang/go-engine/adt/compvar:Data.*
// goverter:extend DataToExecSnap1
var (
	// goverter:map . CompRef
	DataToExecRec func(compexec.ExecRecDS) (compexec.ExecRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(compexec.ExecRec) compexec.ExecRecDS
	// goverter:map . CompRef
	DataToExecSnap2 func(compexec.ExecSnap2DS) (compexec.ExecSnap2, error)
	// goverter:autoMap CompRef
	DataFromExecSnap2 func(compexec.ExecSnap2) compexec.ExecSnap2DS

	DataToRefMap  func(map[uniqsym.ADT]compexec.ExecSnap2DS) (map[uniqsym.ADT]compexec.ExecSnap2, error)
	DataToSnapMap func(map[uniqsym.ADT]compexec.ExecSnap1DS) (map[uniqsym.ADT]compexec.ExecSnap1, error)
)
