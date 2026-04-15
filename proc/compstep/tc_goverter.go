package compstep

import (
	"github.com/orglang/go-sdk/adt/procstep"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/proc/termexp:Msg.*
var (
	MsgFromStepSpec func(StepSpec) procstep.StepSpec
	MsgToStepSpec   func(procstep.StepSpec) (StepSpec, error)
)
