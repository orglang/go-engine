package typeexp

import (
	"github.com/orglang/go-sdk/proc/typeexp"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend Msg.*
var (
	MsgFromExpRefs func([]ExpRef) []typeexp.ExpRef
	MsgToExpRefs   func([]typeexp.ExpRef) ([]ExpRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend data.*
var (
	DataToExpRefs   func([]expRefDS) ([]ExpRef, error)
	DataFromExpRefs func([]ExpRef) []expRefDS
	DataToExpRecs   func([]expRecDS) ([]ExpRec, error)
	DataFromExpRecs func([]ExpRec) []expRecDS
)
