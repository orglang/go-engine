package procstep

import "github.com/orglang/go-sdk/adt/procstep"

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-runtime/adt/identity:Convert.*
// goverter:extend orglang/go-runtime/adt/procdef:Msg.*
var (
	MsgFromStepSpec func(StepSpec) procstep.StepSpec
	MsgToStepSpec   func(procstep.StepSpec) (StepSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-runtime/adt/identity:Convert.*
var (
	DataToSemRecs   func([]StepRecDS) ([]StepRec, error)
	DataFromSemRecs func([]StepRec) ([]StepRecDS, error)
)
