package adapter

import (
	commact "orglang/go-engine/internal/proc/commact/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToExprSpecs   func([]ExprSpec) ([]commact.ExprSpec, error)
	DataFromExprSpecs func([]commact.ExprSpec) ([]ExprSpec, error)
	DataToExprRecs    func([]ExprRec) ([]commact.ExprRec, error)
	DataFromExprRecs  func([]commact.ExprRec) ([]ExprRec, error)
)
