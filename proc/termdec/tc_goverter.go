package termdec

import (
	"github.com/orglang/go-sdk/proc/termdec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/proc/typedef:Msg.*
var (
	MsgToDecSpec    func(termdec.DecSpec) (DecSpec, error)
	MsgFromDecSpec  func(DecSpec) termdec.DecSpec
	MsgToDecSnap    func(termdec.DecSnap) (DecSnap, error)
	MsgFromDecSnap  func(DecSnap) termdec.DecSnap
	MsgFromDecSnaps func([]DecSnap) []termdec.DecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	ViewFromDecSnap func(DecSnap) DecSnapVP
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
	DataToDecRec func(decRecDS) (DecRec, error)
	// goverter:autoMap TermRef
	DataFromDecRec  func(DecRec) (decRecDS, error)
	DataToDecRecs   func([]decRecDS) ([]DecRec, error)
	DataFromDecRecs func([]DecRec) ([]decRecDS, error)
	// goverter:map . TermRef
	DataToDecSnap func(decSnapDS) (DecSnap, error)
	// goverter:autoMap TermRef
	DataFromDecSnap  func(DecSnap) (decSnapDS, error)
	DataToDecSnaps   func([]decSnapDS) ([]DecSnap, error)
	DataFromDecSnaps func([]DecSnap) ([]decSnapDS, error)
)
