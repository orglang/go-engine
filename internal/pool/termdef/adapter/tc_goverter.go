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
	DecSpecFromGenerated func(api.PoolDecSpec) (termdef.DefSpec, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/termvar/adapter:Data.*
var (
	// goverter:map . TermRef | dataToTermRef
	DataToDecRec func(DefRec) (termdef.DefRec, error)
	// goverter:autoMap TermRef
	DataFromDecRec func(termdef.DefRec) (DefRec, error)
	// goverter:ignore TermRN
	dataToTermRef func(DefRec) (termsem.SemRef, error)
)
