package adapter

import (
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
var (
	ViewFromDefSnap func(typedef.DefSnap) DefSnapVP
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/contkey:Convert.*
var (
	// goverter:map . TypeRef
	DataToDefRec  func(DefRec) (typedef.DefRec, error)
	DataToDefRecs func([]DefRec) ([]typedef.DefRec, error)
	// goverter:autoMap TypeRef
	DataFromDefRec  func(typedef.DefRec) (DefRec, error)
	DataFromDefRecs func([]typedef.DefRec) ([]DefRec, error)
)
