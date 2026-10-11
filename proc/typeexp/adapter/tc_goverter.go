package adapter

import (
	sdktypeexp "github.com/orglang/go-sdk/proc/typeexp"

	typeexp "orglang/go-engine/proc/typeexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend Msg.*
var (
	MsgFromExpRefs func([]typeexp.ExpRef) []sdktypeexp.ExpRef
	MsgToExpRefs   func([]sdktypeexp.ExpRef) ([]typeexp.ExpRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend data.*
var (
	DataToExpRefs   func([]typeexp.ExpRefDS) ([]typeexp.ExpRef, error)
	DataFromExpRefs func([]typeexp.ExpRef) []typeexp.ExpRefDS
	DataToExpRecs   func([]typeexp.ExpRecDS) ([]typeexp.ExpRec, error)
	DataFromExpRecs func([]typeexp.ExpRec) []typeexp.ExpRecDS
)
