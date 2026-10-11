package adapter

import (
	proccommturn "orglang/go-engine/proc/commturn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend data.*
var (
	DataToStepRecs   func([]proccommturn.StepRecDS) ([]proccommturn.TurnRec, error)
	DataFromStepRecs func([]proccommturn.TurnRec) ([]proccommturn.StepRecDS, error)
)
