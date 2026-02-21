package xactdef

import (
	"github.com/orglang/go-sdk/adt/xactdef"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/xactexp:Msg.*
var (
	MsgFromDefSpec  func(DefSpec) xactdef.DefSpec
	MsgToDefSpec    func(xactdef.DefSpec) (DefSpec, error)
	MsgFromDefSnap  func(DefSnap) xactdef.DefSnap
	MsgToDefSnap    func(xactdef.DefSnap) (DefSnap, error)
	MsgFromDefSnaps func([]DefSnap) []xactdef.DefSnap
	MsgToDefSnaps   func([]xactdef.DefSnap) ([]DefSnap, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
var (
	// goverter:map . DescRef
	DataToDefRec func(defRecDS) (DefRec, error)
	// goverter:autoMap DescRef
	DataFromDefRec  func(DefRec) (defRecDS, error)
	DataToDefRecs   func([]defRecDS) ([]DefRec, error)
	DataFromDefRecs func([]DefRec) ([]defRecDS, error)
)
