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
	// goverter:ignore StructVars LinearVars
	MsgToExecSnap   func(poolexec.ExecSnap) (ExecCfgSnap, error)
	MsgFromExecSnap func(ExecCfgSnap) poolexec.ExecSnap
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
	// goverter:map . ImplRef
	DataToExecSnap func(execCfgSnapDS) (ExecCfgSnap, error)
	// goverter:autoMap ImplRef
	DataFromExecSnap func(ExecCfgSnap) execCfgSnapDS
	// goverter:ignore ImplRN
	DataToImplRef func(execRecDS) (implsem.SemRef, error)

	DataToRefMap  func(map[uniqsym.ADT]execRecDS) (map[uniqsym.ADT]ExecRec, error)
	DataToSnapMap func(map[uniqsym.ADT]execLiabSnapDS) (map[uniqsym.ADT]ExecLiabSnap, error)
)
