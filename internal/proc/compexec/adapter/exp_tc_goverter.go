package adapter

import (
	compexp "orglang/go-engine/internal/proc/compexec/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToExpSpecs   func([]ExpSpec) ([]compexp.ExpSpec, error)
	DataFromExpSpecs func([]compexp.ExpSpec) ([]ExpSpec, error)
	DataToExpRecs    func([]ExpRec) ([]compexp.ExpRec, error)
	DataFromExpRecs  func([]compexp.ExpRec) ([]ExpRec, error)
)
