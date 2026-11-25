package adapter

import (
	compexp "orglang/go-engine/internal/proc/compexec/exp"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToCompSpecs   func([]CompSpec) ([]compexp.ExpSpec, error)
	DataFromCompSpecs func([]compexp.ExpSpec) ([]CompSpec, error)
	DataToExpRecs     func([]ExpRec) ([]compexp.ExpRec, error)
	DataFromExpRecs   func([]compexp.ExpRec) ([]ExpRec, error)
)
