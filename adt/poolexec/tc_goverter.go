package poolexec

import (
	"github.com/orglang/go-sdk/adt/poolexec"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/uniqsym"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend ConvertRecToRef
var (
	ConvertRecsToRefs func([]ExecRec) []implsem.SemRef
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/implsem:Msg.*
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
	// goverter:map . ImplRef
	// goverter:ignore ChnlBS
	DataToExecRec func(execRecDS) (ExecRec, error)
	// goverter:autoMap ImplRef
	DataFromExecRec func(ExecRec) execRecDS
	// goverter:map . ImplRef
	DataToExecSnap func(execSnapDS) (ExecSnap, error)
	// goverter:autoMap ImplRef
	DataFromExecSnap func(ExecSnap) execSnapDS
	DataToRefMap     func(map[uniqsym.ADT]execRecDS) (map[uniqsym.ADT]ExecRec, error)
)
