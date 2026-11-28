package adapter

import (
	commact "orglang/go-engine/internal/proc/commact/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToExpSpecs   func([]ExprSpec) ([]commact.ExprSpec, error)
	DataFromExpSpecs func([]commact.ExprSpec) ([]ExprSpec, error)
	DataToExpRecs    func([]ExprRec) ([]commact.ExprRec, error)
	DataFromExpRecs  func([]commact.ExprRec) ([]ExprRec, error)
)
