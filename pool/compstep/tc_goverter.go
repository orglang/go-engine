package compstep

import (
	"github.com/orglang/go-sdk/adt/poolstep"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/semterm:Msg.*
// goverter:extend orglang/go-engine/pool/termexp:Msg.*
var (
	MsgToStepSpec   func(poolstep.StepSpec) (StepSpec, error)
	MsgFromStepSpec func(StepSpec) poolstep.StepSpec
)
