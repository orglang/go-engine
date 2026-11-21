package adapter

import (
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToExpSpecs   func([]ExpSpec) ([]termexp.ExpSpec, error)
	DataFromExpSpecs func([]termexp.ExpSpec) ([]ExpSpec, error)
	DataToExpRecs    func([]ExpRec) ([]termexp.ExpRec, error)
	DataFromExpRecs  func([]termexp.ExpRec) ([]ExpRec, error)
)
