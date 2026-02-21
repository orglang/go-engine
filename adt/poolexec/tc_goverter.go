package poolexec

import (
	"github.com/orglang/go-sdk/adt/poolexec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	MsgToExecSpec   func(poolexec.ExecSpec) (ExecSpec, error)
	MsgFromExecSpec func(ExecSpec) poolexec.ExecSpec
	MsgToExecSnap   func(poolexec.ExecSnap) (ExecSnap, error)
	MsgFromExecSnap func(ExecSnap) poolexec.ExecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
var (
	// goverter:map . ExecRef
	DataToExecRec func(execRecDS) (ExecRec, error)
	// goverter:autoMap ExecRef
	DataFromExecRec func(ExecRec) execRecDS
	// goverter:map . ExecRef
	DataToLiab func(liabDS) (Liab, error)
	// goverter:autoMap ExecRef
	DataFromLiab func(Liab) liabDS
	// goverter:map . ExecRef
	DataToExecSnap func(execSnapDS) (ExecSnap, error)
	// goverter:autoMap ExecRef
	DataFromExecSnap func(ExecSnap) execSnapDS
)
