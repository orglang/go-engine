package adapter

import (
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/contkey:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
var (
	// goverter:map . DescRef
	// goverter:ignore TypeExp
	DataToDefRec  func(DefRec) (typedef.DefRec, error)
	DataToDefRecs func([]DefRec) ([]typedef.DefRec, error)
	// goverter:autoMap DescRef
	DataFromDefRec  func(typedef.DefRec) (DefRec, error)
	DataFromDefRecs func([]typedef.DefRec) ([]DefRec, error)
)
