package typedef

import (
	"github.com/orglang/go-sdk/adt/typedef"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/typeexp:Msg.*
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
// goverter:extend orglang/go-engine/adt/typeexp:Msg.*
var (
	ViewFromDefSnap func(DefSnap) DefSnapVP
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
// goverter:extend orglang/go-engine/adt/typeexp:Data.*
var (
	// goverter:map . DescRef
	DataToDefRec func(defRecDS) (DefRec, error)
	// goverter:autoMap DescRef
	DataFromDefRec  func(DefRec) (defRecDS, error)
	DataToDefRecs   func([]defRecDS) ([]DefRec, error)
	DataFromDefRecs func([]DefRec) ([]defRecDS, error)
)
