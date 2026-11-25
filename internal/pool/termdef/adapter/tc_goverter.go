package adapter

import (
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/app/web/api"

	termdef "orglang/go-engine/internal/pool/termdef/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/symbol:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
var (
	// goverter:map ConnPh ChnlPH
	// goverter:map TypeQn TypeQN
	ConnVarFromGenerated func(api.AnyConnVarSpec) (termvar.VarSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend ConnVarFromGenerated
var (
	// goverter:map DescQn TermQN
	DecSpecFromGenerated func(api.PoolDecSpec) (termdef.DecSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/contkey:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/termvar/adapter:Data.*
var (
	// goverter:map . TermRef | dataToTermRef
	// goverter:ignore ExpCK
	DataToDecRec func(DecRec) (termdef.DecRec, error)
	// goverter:autoMap TermRef
	// goverter:ignore ExpCK
	DataFromDecRec func(termdef.DecRec) (DecRec, error)
	dataToTermRef  func(DecRec) (termsem.SemRef, error)
)
