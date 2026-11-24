package adapter

import (
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToExpSpecs   func([]ExpSpec) ([]termexp.CompSpec, error)
	DataFromExpSpecs func([]termexp.CompSpec) ([]ExpSpec, error)
	DataToExpRecs    func([]ExpRec) ([]termexp.CompRec, error)
	DataFromExpRecs  func([]termexp.CompRec) ([]ExpRec, error)
)
