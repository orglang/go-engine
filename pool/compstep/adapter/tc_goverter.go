package adapter

import (
	"github.com/orglang/go-sdk/pool/compstep"

	poolcompstep "orglang/go-engine/pool/compstep/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/compsem:Msg.*
// goverter:extend orglang/go-engine/pool/termexp/adapter:Msg.*
var (
	MsgToStepSpec   func(compstep.StepSpec) (poolcompstep.StepSpec, error)
	MsgFromStepSpec func(poolcompstep.StepSpec) compstep.StepSpec
)
