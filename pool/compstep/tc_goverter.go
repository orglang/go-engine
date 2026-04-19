package compstep

import (
	"github.com/orglang/go-sdk/pool/compstep"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/compsem:Msg.*
// goverter:extend orglang/go-engine/pool/termexp:Msg.*
var (
	MsgToStepSpec   func(compstep.StepSpec) (StepSpec, error)
	MsgFromStepSpec func(StepSpec) compstep.StepSpec
)
