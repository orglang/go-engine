package termdef

import (
	"orglang/go-engine/adt/termsem"

	"github.com/orglang/go-sdk/pool/termdec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/pool/typedef:Msg.*
var (
	MsgToDecSpec   func(termdec.DecSpec) (DefSpec, error)
	MsgFromDecSpec func(DefSpec) termdec.DecSpec
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/termvar:Data.*
var (
	// goverter:map . TermRef | dataToTermRef
	DataToDecRec func(defRecDS) (DefRec, error)
	// goverter:autoMap TermRef
	DataFromDecRec func(DefRec) (defRecDS, error)
	// goverter:ignore TermRN
	dataToTermRef func(defRecDS) (termsem.SemRef, error)
)
