package adapter

import (
	sdktypeexp "github.com/orglang/go-sdk/pool/typeexp"

	typeexp "orglang/go-engine/pool/typeexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend msg.*
var (
	MsgFromExpRefs func([]typeexp.ExpRef) []sdktypeexp.ExpRef
	MsgToExpRefs   func([]sdktypeexp.ExpRef) ([]typeexp.ExpRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend data.*
var (
	DataToExpRefs   func([]typeexp.ExpRefDS) ([]typeexp.ExpRef, error)
	DataFromExpRefs func([]typeexp.ExpRef) []typeexp.ExpRefDS
	DataToExpRecs   func([]typeexp.ExpRecDS) ([]typeexp.ExpRec, error)
	DataFromExpRecs func([]typeexp.ExpRec) []typeexp.ExpRecDS
)
