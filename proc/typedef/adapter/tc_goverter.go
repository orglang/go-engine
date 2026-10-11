package adapter

import (
	"github.com/orglang/go-sdk/proc/typedef"

	proctypedef "orglang/go-engine/proc/typedef/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/proc/typeexp/adapter:Msg.*
var (
	MsgFromDefSpec  func(proctypedef.DefSpec) typedef.DefSpec
	MsgToDefSpec    func(typedef.DefSpec) (proctypedef.DefSpec, error)
	MsgFromDefSnap  func(proctypedef.DefSnap) typedef.DefSnap
	MsgToDefSnap    func(typedef.DefSnap) (proctypedef.DefSnap, error)
	MsgFromDefSnaps func([]proctypedef.DefSnap) []typedef.DefSnap
	MsgToDefSnaps   func([]typedef.DefSnap) ([]proctypedef.DefSnap, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/proc/typeexp/adapter:Msg.*
var (
	ViewFromDefSnap func(proctypedef.DefSnap) DefSnapVP
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
// goverter:extend orglang/go-engine/proc/typeexp/adapter:Data.*
var (
	// goverter:map . TypeRef
	DataToDefRec  func(proctypedef.DefRecDS) (proctypedef.DefRec, error)
	DataToDefRecs func([]proctypedef.DefRecDS) ([]proctypedef.DefRec, error)
	// goverter:autoMap TypeRef
	DataFromDefRec  func(proctypedef.DefRec) (proctypedef.DefRecDS, error)
	DataFromDefRecs func([]proctypedef.DefRec) ([]proctypedef.DefRecDS, error)
)
