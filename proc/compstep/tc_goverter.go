package compstep

import "github.com/orglang/go-sdk/proc/compstep"

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/proc/termexp:Msg.*
var (
	MsgFromStepSpec func(StepSpec) compstep.StepSpec
	MsgToStepSpec   func(compstep.StepSpec) (StepSpec, error)
)
