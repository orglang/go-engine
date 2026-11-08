package adapter

import (
	"orglang/go-engine/internal/adt/termsem"

	"github.com/orglang/go-sdk/pool/termdec"

	termdef "orglang/go-engine/internal/pool/termdef/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/pool/typedef/adapter:Msg.*
var (
	MsgToDecSpec   func(termdec.DecSpec) (termdef.DefSpec, error)
	MsgFromDecSpec func(termdef.DefSpec) termdec.DecSpec
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/termvar:Data.*
var (
	// goverter:map . TermRef | dataToTermRef
	DataToDecRec func(termdef.DefRecDS) (termdef.DefRec, error)
	// goverter:autoMap TermRef
	DataFromDecRec func(termdef.DefRec) (termdef.DefRecDS, error)
	// goverter:ignore TermRN
	dataToTermRef func(termdef.DefRecDS) (termsem.SemRef, error)
)
