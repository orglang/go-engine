package poolexec

import (
	"github.com/orglang/go-sdk/adt/poolexec"

	"orglang/go-engine/adt/semterm"
	"orglang/go-engine/adt/uniqsym"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend ConvertRecToRef
var (
	ConvertRecsToRefs func([]ExecRec) []semterm.TermRef
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/semterm:Msg.*
var (
	MsgToExecSpec   func(poolexec.ExecSpec) (ExecSpec, error)
	MsgFromExecSpec func(ExecSpec) poolexec.ExecSpec
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/implvar:Convert.*
// goverter:extend orglang/go-engine/adt/implvar:Data.*
// goverter:extend DataToExecLiabSnap
var (
	// goverter:map . ImplRef | DataToImplRef
	DataToExecRec func(execRecDS) (ExecRec, error)
	// goverter:autoMap ImplRef
	DataFromExecRec func(ExecRec) execRecDS
	// goverter:ignore ImplRN
	DataToImplRef func(execRecDS) (semterm.TermRef, error)

	DataToRefMap  func(map[uniqsym.ADT]execRecDS) (map[uniqsym.ADT]ExecRec, error)
	DataToSnapMap func(map[uniqsym.ADT]execSnapDS) (map[uniqsym.ADT]ExecSnap, error)
)
