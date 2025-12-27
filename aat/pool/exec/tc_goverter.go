package exec

import (
	procexec "orglang/orglang/aat/proc/exec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/orglang/avt/id:Convert.*
var (
	ConvertRecToRef func(PoolRec) PoolRef
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/orglang/avt/id:Convert.*
// goverter:extend orglang/orglang/aat/proc/def:Msg.*
var (
	MsgToPoolSpec   func(PoolSpecME) (PoolSpec, error)
	MsgFromPoolSpec func(PoolSpec) PoolSpecME
	MsgToPoolRef    func(PoolRefME) (PoolRef, error)
	MsgFromPoolRef  func(PoolRef) PoolRefME
	MsgToPoolSnap   func(PoolSnapME) (PoolSnap, error)
	MsgFromPoolSnap func(PoolSnap) PoolSnapME
	MsgFromStepSpec func(StepSpec) StepSpecME
	MsgToStepSpec   func(StepSpecME) (StepSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/orglang/avt/id:Convert.*
var (
	DataToPoolRef    func(poolRefDS) (PoolRef, error)
	DataFromPoolRef  func(PoolRef) poolRefDS
	DataToPoolRefs   func([]poolRefDS) ([]PoolRef, error)
	DataFromPoolRefs func([]PoolRef) []poolRefDS
	DataToPoolRec    func(poolRecDS) (PoolRec, error)
	DataFromPoolRec  func(PoolRec) poolRecDS
	DataToLiab       func(liabDS) (procexec.Liab, error)
	DataFromLiab     func(procexec.Liab) liabDS
	DataToPoolSnap   func(poolSnapDS) (PoolSnap, error)
	DataFromPoolSnap func(PoolSnap) poolSnapDS
	DataToEPs        func([]epDS) ([]procexec.EP, error)
)
