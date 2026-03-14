package poolstep

import (
	"github.com/orglang/go-sdk/adt/poolstep"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/implsem:Msg.*
// goverter:extend orglang/go-engine/adt/poolexp:Msg.*
var (
	MsgToCommSpec   func(poolstep.StepSpec) (StepSpec, error)
	MsgFromCommSpec func(StepSpec) poolstep.StepSpec
)
