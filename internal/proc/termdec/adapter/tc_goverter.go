package adapter

import (
	"orglang/go-engine/internal/adt/termvar"
	webapi "orglang/go-engine/internal/app/web/api"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/adt/contkey:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/termvar/adapter:Data.*
var (
	// goverter:map . DescRef
	// goverter:ignore TermQN
	DataToTermDecRec func(DecRec) (termdec.DecRec, error)
	// goverter:autoMap DescRef
	DataFromTermDecRec  func(termdec.DecRec) (DecRec, error)
	DataToTermDecRecs   func([]DecRec) ([]termdec.DecRec, error)
	DataFromTermDecRecs func([]termdec.DecRec) ([]DecRec, error)
	// goverter:map . DescRef
	DataToTermDecSnap func(DecSnap) (termdec.DecSnap, error)
	// goverter:autoMap DescRef
	DataFromTermDecSnap  func(termdec.DecSnap) (DecSnap, error)
	DataToTermDecSnaps   func([]DecSnap) ([]termdec.DecSnap, error)
	DataFromTermDecSnaps func([]termdec.DecSnap) ([]DecSnap, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/adt/symbol:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
var (
	// goverter:map CommPh ChnlPH
	// goverter:map TypeQn TypeQN
	ConnVarFromGenerated func(webapi.AnyConnVarSpec) (termvar.VarSpec, error)
	// goverter:map TermQn TermQN
	DecSpecFromGenerated func(webapi.ProcDecSpec) (termdec.DecSpec, error)
	// goverter:map TermID TermId
	// goverter:map DescRef.DescRN DescRn
	DecSnapToGenerated func(termdec.DecSnap) webapi.AnyTermRef
)
