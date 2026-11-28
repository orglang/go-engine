package adapter

import (
	commact "orglang/go-engine/internal/proc/commact/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToExpSpecs   func([]ExpSpec) ([]commact.ExpSpec, error)
	DataFromExpSpecs func([]commact.ExpSpec) ([]ExpSpec, error)
	DataToExpRecs    func([]ExpRec) ([]commact.ExpRec, error)
	DataFromExpRecs  func([]commact.ExpRec) ([]ExpRec, error)
)
