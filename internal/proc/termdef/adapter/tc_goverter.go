package adapter

import (
	"orglang/go-engine/internal/adt/termvar"
	webapi "orglang/go-engine/internal/app/web/api"
	termdef "orglang/go-engine/internal/proc/termdef/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/contkey:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
var (
	ViewFromDecSnap func(termdef.DecSnap) DecSnapVP
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/adt/contkey:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/termvar/adapter:Data.*
var (
	// goverter:map . TermRef
	// goverter:ignore TermQN
	// goverter:ignore ExpCK
	DataToDecRec func(DecRec) (termdef.DecRec, error)
	// goverter:autoMap TermRef
	// goverter:ignore ExpCK
	DataFromDecRec  func(termdef.DecRec) (DecRec, error)
	DataToDecRecs   func([]DecRec) ([]termdef.DecRec, error)
	DataFromDecRecs func([]termdef.DecRec) ([]DecRec, error)
	// goverter:map . TermRef
	DataToDecSnap func(DecSnap) (termdef.DecSnap, error)
	// goverter:autoMap TermRef
	DataFromDecSnap  func(termdef.DecSnap) (DecSnap, error)
	DataToDecSnaps   func([]DecSnap) ([]termdef.DecSnap, error)
	DataFromDecSnaps func([]termdef.DecSnap) ([]DecSnap, error)
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
	DecSpecFromGenerated func(webapi.ProcDecSpec) (termdef.DecSpec, error)

	// goverter:map TermRef.TermID TermId
	// goverter:map TermRef.DescRN DescRn
	DecSnapToGenerated func(termdef.DecSnap) webapi.AnyTermRef
)
