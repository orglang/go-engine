package compexec

import (
	"github.com/orglang/go-sdk/pool/compexec"

	"orglang/go-engine/adt/uniqsym"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/compsem:Msg.*
var (
	MsgToExecSpec   func(compexec.ExecSpec) (ExecSpec, error)
	MsgFromExecSpec func(ExecSpec) compexec.ExecSpec
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/compvar:Convert.*
// goverter:extend orglang/go-engine/adt/compvar:Data.*
// goverter:extend DataToExecSnap1
var (
	// goverter:map . CompRef
	DataToExecRec func(execRec) (ExecRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(ExecRec) execRec
	// goverter:map . CompRef
	DataToExecSnap2 func(execSnap2) (ExecSnap2, error)
	// goverter:autoMap CompRef
	DataFromExecSnap2 func(ExecSnap2) execSnap2

	DataToRefMap  func(map[uniqsym.ADT]execSnap2) (map[uniqsym.ADT]ExecSnap2, error)
	DataToSnapMap func(map[uniqsym.ADT]execSnap1) (map[uniqsym.ADT]ExecSnap1, error)
)
