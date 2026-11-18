package adapter

import (
	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend data.*
var (
	DataToExpRefs   func([]typeexp.ExpRefDS) ([]typeexp.ExpRef, error)
	DataFromExpRefs func([]typeexp.ExpRef) []typeexp.ExpRefDS
	DataToExpRecs   func([]typeexp.ExpRecDS) ([]typeexp.ExpRec, error)
	DataFromExpRecs func([]typeexp.ExpRec) []typeexp.ExpRecDS
)
