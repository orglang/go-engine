package adapter

import (
	sdkcompstep "github.com/orglang/go-sdk/proc/compstep"

	compstep "orglang/go-engine/proc/compstep/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/proc/termexp/adapter:Msg.*
var (
	MsgFromStepSpec func(compstep.StepSpec) sdkcompstep.StepSpec
	MsgToStepSpec   func(sdkcompstep.StepSpec) (compstep.StepSpec, error)
)
