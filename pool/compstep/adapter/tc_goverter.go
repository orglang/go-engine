package adapter

import (
	sdkcompstep "github.com/orglang/go-sdk/pool/compstep"

	compstep "orglang/go-engine/pool/compstep/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/compsem:Msg.*
// goverter:extend orglang/go-engine/pool/termexp/adapter:Msg.*
var (
	MsgToStepSpec   func(sdkcompstep.StepSpec) (compstep.StepSpec, error)
	MsgFromStepSpec func(compstep.StepSpec) sdkcompstep.StepSpec
)
