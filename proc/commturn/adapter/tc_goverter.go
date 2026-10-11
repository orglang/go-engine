package adapter

import (
	commturn "orglang/go-engine/proc/commturn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend data.*
var (
	DataToStepRecs   func([]commturn.StepRecDS) ([]commturn.TurnRec, error)
	DataFromStepRecs func([]commturn.TurnRec) ([]commturn.StepRecDS, error)
)
