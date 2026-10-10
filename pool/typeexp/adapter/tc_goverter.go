package adapter

import (
	"github.com/orglang/go-sdk/pool/typeexp"

	pooltypeexp "orglang/go-engine/pool/typeexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend msg.*
var (
	MsgFromExpRefs func([]pooltypeexp.ExpRef) []typeexp.ExpRef
	MsgToExpRefs   func([]typeexp.ExpRef) ([]pooltypeexp.ExpRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend data.*
var (
	DataToExpRefs   func([]pooltypeexp.ExpRefDS) ([]pooltypeexp.ExpRef, error)
	DataFromExpRefs func([]pooltypeexp.ExpRef) []pooltypeexp.ExpRefDS
	DataToExpRecs   func([]pooltypeexp.ExpRecDS) ([]pooltypeexp.ExpRec, error)
	DataFromExpRecs func([]pooltypeexp.ExpRec) []pooltypeexp.ExpRecDS
)
