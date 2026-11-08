package adapter

import (
	sdktermdec "github.com/orglang/go-sdk/proc/termdec"

	termdec "orglang/go-engine/proc/termdec/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/proc/typedef/adapter:Msg.*
var (
	MsgToDecSpec    func(sdktermdec.DecSpec) (termdec.DecSpec, error)
	MsgFromDecSpec  func(termdec.DecSpec) sdktermdec.DecSpec
	MsgToDecSnap    func(sdktermdec.DecSnap) (termdec.DecSnap, error)
	MsgFromDecSnap  func(termdec.DecSnap) sdktermdec.DecSnap
	MsgFromDecSnaps func([]termdec.DecSnap) []sdktermdec.DecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	ViewFromDecSnap func(termdec.DecSnap) DecSnapVP
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
// goverter:extend orglang/go-engine/adt/termvar:Data.*
var (
	// goverter:map . TermRef
	// goverter:ignore TermQN
	DataToDecRec func(termdec.DecRecDS) (termdec.DecRec, error)
	// goverter:autoMap TermRef
	DataFromDecRec  func(termdec.DecRec) (termdec.DecRecDS, error)
	DataToDecRecs   func([]termdec.DecRecDS) ([]termdec.DecRec, error)
	DataFromDecRecs func([]termdec.DecRec) ([]termdec.DecRecDS, error)
	// goverter:map . TermRef
	DataToDecSnap func(termdec.DecSnapDS) (termdec.DecSnap, error)
	// goverter:autoMap TermRef
	DataFromDecSnap  func(termdec.DecSnap) (termdec.DecSnapDS, error)
	DataToDecSnaps   func([]termdec.DecSnapDS) ([]termdec.DecSnap, error)
	DataFromDecSnaps func([]termdec.DecSnap) ([]termdec.DecSnapDS, error)
)
