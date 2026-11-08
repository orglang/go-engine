package adapter

import (
	sdkcompstep "github.com/orglang/go-sdk/pool/compstep"

	compstep "orglang/go-engine/internal/pool/compstep/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/compsem:Msg.*
// goverter:extend orglang/go-engine/internal/pool/termexp/adapter:Msg.*
var (
	MsgToStepSpec   func(sdkcompstep.StepSpec) (compstep.StepSpec, error)
	MsgFromStepSpec func(compstep.StepSpec) sdkcompstep.StepSpec
)
