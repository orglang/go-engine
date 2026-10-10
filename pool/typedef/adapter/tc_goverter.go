package adapter

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/pool/typedef"

	pooltypedef "orglang/go-engine/pool/typedef/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/pool/typeexp/adapter:Msg.*
var (
	MsgFromDefSpec  func(pooltypedef.DefSpec) typedef.DefSpec
	MsgToDefSpec    func(typedef.DefSpec) (pooltypedef.DefSpec, error)
	MsgFromDefSnap  func(pooltypedef.DefSnap) typedef.DefSnap
	MsgToDefSnap    func(typedef.DefSnap) (pooltypedef.DefSnap, error)
	MsgFromDefSnaps func([]pooltypedef.DefSnap) []typedef.DefSnap
	MsgToDefSnaps   func([]typedef.DefSnap) ([]pooltypedef.DefSnap, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
var (
	// goverter:map . TypeRef
	DataToDefRec    func(pooltypedef.DefRecDS) (pooltypedef.DefRec, error)
	DataToDefRecs   func([]pooltypedef.DefRecDS) ([]pooltypedef.DefRec, error)
	DataToDefRecMap func(map[uniqsym.ADT]pooltypedef.DefRecDS) (map[uniqsym.ADT]pooltypedef.DefRec, error)
	// goverter:autoMap TypeRef
	DataFromDefRec  func(pooltypedef.DefRec) (pooltypedef.DefRecDS, error)
	DataFromDefRecs func([]pooltypedef.DefRec) ([]pooltypedef.DefRecDS, error)
)
