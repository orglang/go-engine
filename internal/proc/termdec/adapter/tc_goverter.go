package adapter

import (
	"orglang/go-engine/internal/adt/termvar"
	webapi "orglang/go-engine/internal/app/web/api"
	termdec "orglang/go-engine/internal/proc/termdec/core"
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
// goverter:extend orglang/go-engine/internal/adt/termvar/adapter:Data.*
var (
	// goverter:map . TermRef
	// goverter:ignore TermQN
	DataToDecRec func(DecRecDS) (termdec.DecRec, error)
	// goverter:autoMap TermRef
	DataFromDecRec  func(termdec.DecRec) (DecRecDS, error)
	DataToDecRecs   func([]DecRecDS) ([]termdec.DecRec, error)
	DataFromDecRecs func([]termdec.DecRec) ([]DecRecDS, error)
	// goverter:map . TermRef
	DataToDecSnap func(DecSnapDS) (termdec.DecSnap, error)
	// goverter:autoMap TermRef
	DataFromDecSnap  func(termdec.DecSnap) (DecSnapDS, error)
	DataToDecSnaps   func([]DecSnapDS) ([]termdec.DecSnap, error)
	DataFromDecSnaps func([]termdec.DecSnap) ([]DecSnapDS, error)
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
