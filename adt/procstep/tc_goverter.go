package procstep

import (
	"github.com/orglang/go-sdk/adt/procstep"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/procexp:Msg.*
var (
	MsgFromStepSpec func(CommSpec) procstep.StepSpec
	MsgToStepSpec   func(procstep.StepSpec) (CommSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend data.*
var (
	DataToStepRecs   func([]StepRecDS) ([]CommRec, error)
	DataFromStepRecs func([]CommRec) ([]StepRecDS, error)
)
