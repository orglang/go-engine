package adapter

import (
	"orglang/go-engine/adt/uniqsym"

	sdktypedef "github.com/orglang/go-sdk/pool/typedef"

	typedef "orglang/go-engine/pool/typedef/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/pool/typeexp/adapter:Msg.*
var (
	MsgFromDefSpec  func(typedef.DefSpec) sdktypedef.DefSpec
	MsgToDefSpec    func(sdktypedef.DefSpec) (typedef.DefSpec, error)
	MsgFromDefSnap  func(typedef.DefSnap) sdktypedef.DefSnap
	MsgToDefSnap    func(sdktypedef.DefSnap) (typedef.DefSnap, error)
	MsgFromDefSnaps func([]typedef.DefSnap) []sdktypedef.DefSnap
	MsgToDefSnaps   func([]sdktypedef.DefSnap) ([]typedef.DefSnap, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
var (
	// goverter:map . TypeRef
	DataToDefRec    func(typedef.DefRecDS) (typedef.DefRec, error)
	DataToDefRecs   func([]typedef.DefRecDS) ([]typedef.DefRec, error)
	DataToDefRecMap func(map[uniqsym.ADT]typedef.DefRecDS) (map[uniqsym.ADT]typedef.DefRec, error)
	// goverter:autoMap TypeRef
	DataFromDefRec  func(typedef.DefRec) (typedef.DefRecDS, error)
	DataFromDefRecs func([]typedef.DefRec) ([]typedef.DefRecDS, error)
)
