package adapter

import (
	connact "orglang/go-engine/internal/proc/connact/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToExpSpecs   func([]ExpSpec) ([]connact.ExpSpec, error)
	DataFromExpSpecs func([]connact.ExpSpec) ([]ExpSpec, error)
	DataToExpRecs    func([]ExpRec) ([]connact.ExpRec, error)
	DataFromExpRecs  func([]connact.ExpRec) ([]ExpRec, error)
)
