package adapter

import (
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToCompSpecs   func([]CompSpec) ([]termexp.CompSpec, error)
	DataFromCompSpecs func([]termexp.CompSpec) ([]CompSpec, error)
	DataToExpRecs     func([]ExpRec) ([]termexp.CompRec, error)
	DataFromExpRecs   func([]termexp.CompRec) ([]ExpRec, error)
)
