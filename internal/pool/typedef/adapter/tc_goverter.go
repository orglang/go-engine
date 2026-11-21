package adapter

import (
	"orglang/go-engine/internal/adt/uniqsym"

	typedef "orglang/go-engine/internal/pool/typedef/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/contkey:Convert.*
var (
	// goverter:map . TypeRef
	DataToDefRec    func(DefRec) (typedef.DefRec, error)
	DataToDefRecs   func([]DefRec) ([]typedef.DefRec, error)
	DataToDefRecMap func(map[uniqsym.ADT]DefRec) (map[uniqsym.ADT]typedef.DefRec, error)
	// goverter:autoMap TypeRef
	DataFromDefRec  func(typedef.DefRec) (DefRec, error)
	DataFromDefRecs func([]typedef.DefRec) ([]DefRec, error)
)


// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/contkey:Convert.*
var (
	// goverter:map . TypeRef
	DataToDefRecEnt func(DefRecDS) (typedef.DefRec, error)
	// goverter:autoMap TypeRef
	DataFromDefRecEnt func(typedef.DefRec) (DefRecDS, error)
	DataToDefRecMapEnt func(map[uniqsym.ADT]DefRecDS) (map[uniqsym.ADT]typedef.DefRec, error)
)
