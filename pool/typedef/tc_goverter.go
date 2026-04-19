package typedef

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/pool/typedef"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/pool/typeexp:Msg.*
var (
	MsgFromDefSpec  func(DefSpec) typedef.DefSpec
	MsgToDefSpec    func(typedef.DefSpec) (DefSpec, error)
	MsgFromDefSnap  func(DefSnap) typedef.DefSnap
	MsgToDefSnap    func(typedef.DefSnap) (DefSnap, error)
	MsgFromDefSnaps func([]DefSnap) []typedef.DefSnap
	MsgToDefSnaps   func([]typedef.DefSnap) ([]DefSnap, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
var (
	// goverter:map . TypeRef
	// goverter:ignore TypeQN
	DataToDefRec    func(defRecDS) (DefRec, error)
	DataToDefRecs   func([]defRecDS) ([]DefRec, error)
	DataToDefRecMap func(map[uniqsym.ADT]defRecDS) (map[uniqsym.ADT]DefRec, error)
	// goverter:autoMap TypeRef
	DataFromDefRec  func(DefRec) (defRecDS, error)
	DataFromDefRecs func([]DefRec) ([]defRecDS, error)
)
