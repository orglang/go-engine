package procexec

import "github.com/orglang/go-sdk/adt/procexec"

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-runtime/adt/identity:Convert.*
// goverter:extend orglang/go-runtime/adt/procdef:Msg.*
var (
	MsgFromExecSpec func(ExecSpec) procexec.ExecSpec
	MsgToExecSpec   func(procexec.ExecSpec) (ExecSpec, error)
	MsgToExecRef    func(procexec.ExecRef) (ExecRef, error)
	MsgFromExecRef  func(ExecRef) procexec.ExecRef
	MsgToExecSnap   func(procexec.ExecSnap) (ExecSnap, error)
	MsgFromExecSnap func(ExecSnap) procexec.ExecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-runtime/adt/identity:Convert.*
// goverter:extend orglang/go-runtime/adt/revnum:Convert.*
// goverter:extend orglang/go-runtime/adt/procdef:Data.*
// goverter:extend data.*
var (
	DataFromMod  func(Mod) (modDS, error)
	DataFromBnd  func(Bnd) bindDS
	DataToLiab   func(liabDS) (Liab, error)
	DataFromLiab func(Liab) liabDS
	DataToEPs    func([]epDS) ([]EP, error)
)
