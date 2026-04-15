package typedef

import (
	"orglang/go-engine/adt/semtype"
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/xactdef"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/pool/typeexp:Msg.*
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
	// goverter:map . DescRef | DataToSemRef
	DataToDefRec    func(defRecDS) (DefRec, error)
	DataToDefRecs   func([]defRecDS) ([]DefRec, error)
	DataToDefRecMap func(map[uniqsym.ADT]defRecDS) (map[uniqsym.ADT]DefRec, error)
	// goverter:autoMap DescRef
	DataFromDefRec  func(DefRec) (defRecDS, error)
	DataFromDefRecs func([]DefRec) ([]defRecDS, error)
	// goverter:ignore DescRN
	DataToSemRef func(defRecDS) (semtype.TypeRef, error)
)
