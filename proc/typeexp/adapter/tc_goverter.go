package adapter

import (
	"github.com/orglang/go-sdk/proc/typeexp"

	proctypeexp "orglang/go-engine/proc/typeexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend Msg.*
var (
	MsgFromExpRefs func([]proctypeexp.ExpRef) []typeexp.ExpRef
	MsgToExpRefs   func([]typeexp.ExpRef) ([]proctypeexp.ExpRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend data.*
var (
	DataToExpRefs   func([]proctypeexp.ExpRefDS) ([]proctypeexp.ExpRef, error)
	DataFromExpRefs func([]proctypeexp.ExpRef) []proctypeexp.ExpRefDS
	DataToExpRecs   func([]proctypeexp.ExpRecDS) ([]proctypeexp.ExpRec, error)
	DataFromExpRecs func([]proctypeexp.ExpRec) []proctypeexp.ExpRecDS
)
