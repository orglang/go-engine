package adapter

import (
	sdktermdec "github.com/orglang/go-sdk/proc/termdec"

	"orglang/go-engine/internal/adt/termvar"
	webapi "orglang/go-engine/internal/app/web/api"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/proc/typedef/adapter:Msg.*
// goverter:extend MsgFromTermVarRec
// goverter:extend MsgToTermVarRec
var (
	MsgToDecSpec    func(sdktermdec.DecSpec) (termdec.DecSpec, error)
	MsgFromDecSpec  func(termdec.DecSpec) sdktermdec.DecSpec
	MsgToDecSnap    func(sdktermdec.DecSnap) (termdec.DecSnap, error)
	MsgFromDecSnap  func(termdec.DecSnap) sdktermdec.DecSnap
	MsgFromDecSnaps func([]termdec.DecSnap) []sdktermdec.DecSnap
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
var (
	ViewFromDecSnap func(termdec.DecSnap) DecSnapVP
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/adt/contkey:Convert.*
// goverter:extend orglang/go-engine/internal/adt/termvar:Data.*
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

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/adt/symbol:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
var (
	// goverter:map ConnPh ChnlPH
	// goverter:map TypeQn TypeQN
	ConnVarFromGenerated func(webapi.AnyConnVarSpec) (termvar.VarSpec, error)

	// goverter:map TermQn TermQN
	DecSpecFromGenerated func(webapi.ProcDecSpec) (termdec.DecSpec, error)

	// goverter:map TermRef.TermID TermId
	// goverter:map TermRef.TermRN TermRn
	DecSnapToGenerated func(termdec.DecSnap) webapi.AnyTermRef
)
