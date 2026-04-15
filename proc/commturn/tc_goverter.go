package commturn

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend data.*
var (
	DataToStepRecs   func([]StepRecDS) ([]TurnRec, error)
	DataFromStepRecs func([]TurnRec) ([]StepRecDS, error)
)
