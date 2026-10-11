package adapter

import (
	"github.com/orglang/go-sdk/proc/compstep"

	proccompstep "orglang/go-engine/proc/compstep/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/proc/termexp/adapter:Msg.*
var (
	MsgFromStepSpec func(proccompstep.StepSpec) compstep.StepSpec
	MsgToStepSpec   func(compstep.StepSpec) (proccompstep.StepSpec, error)
)
