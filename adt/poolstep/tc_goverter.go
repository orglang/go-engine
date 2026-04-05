package poolstep

import (
	"orglang/go-engine/adt/commsem"

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

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	DataToSemRef func(StepRecDS) (commsem.SemRef, error)
)
