package xactdef

import (
	"github.com/orglang/go-sdk/adt/xactdef"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/uniqref:Msg.*
// goverter:extend orglang/go-engine/adt/xactexp:Msg.*
var (
	MsgFromDefSpec  func(DefSpec) xactdef.DefSpec
	MsgToDefSpec    func(xactdef.DefSpec) (DefSpec, error)
	MsgFromDefRef   func(DefRef) xactdef.DefRef
	MsgToDefRef     func(xactdef.DefRef) (DefRef, error)
	MsgFromDefRefs  func([]DefRef) []xactdef.DefRef
	MsgToDefRefs    func([]xactdef.DefRef) ([]DefRef, error)
	MsgFromDefSnap  func(DefSnap) xactdef.DefSnap
	MsgToDefSnap    func(xactdef.DefSnap) (DefSnap, error)
	MsgFromDefSnaps func([]DefSnap) []xactdef.DefSnap
	MsgToDefSnaps   func([]xactdef.DefSnap) ([]DefSnap, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqref:Data.*
var (
	DataToDefRef    func(defRefDS) (DefRef, error)
	DataFromDefRef  func(DefRef) (defRefDS, error)
	DataToDefRefs   func([]defRefDS) ([]DefRef, error)
	DataFromDefRefs func([]DefRef) ([]defRefDS, error)
	// goverter:map . DefRef
	DataToDefRec func(defRecDS) (DefRec, error)
	// goverter:autoMap DefRef
	DataFromDefRec  func(DefRec) (defRecDS, error)
	DataToDefRecs   func([]defRecDS) ([]DefRec, error)
	DataFromDefRecs func([]DefRec) ([]defRecDS, error)
)
