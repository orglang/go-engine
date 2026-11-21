package adapter

import (
	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend data.*
var (
	DataToExpRefs   func([]ExpRef) ([]typeexp.ExpRef, error)
	DataFromExpRefs func([]typeexp.ExpRef) []ExpRef
	DataToExpRecs   func([]ExpRec) ([]typeexp.ExpRec, error)
	DataFromExpRecs func([]typeexp.ExpRec) []ExpRec
)
