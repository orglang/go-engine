package poolstep

import (
	"github.com/orglang/go-sdk/adt/poolstep"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/implsem:Msg.*
// goverter:extend orglang/go-engine/adt/poolexp:Msg.*
var (
	MsgToStepSpec   func(poolstep.StepSpec) (StepSpec, error)
	MsgFromStepSpec func(StepSpec) poolstep.StepSpec
)
