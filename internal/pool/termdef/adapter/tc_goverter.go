package adapter

import (
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/pool/termdec/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/contkey:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/termvar/adapter:Data.*
var (
	// goverter:map . TermRef | dataToTermRef
	DataToTermDecRec func(DecRec) (core.DecRec, error)
	// goverter:autoMap TermRef
	// goverter:ignore ExpCK
	DataFromTermDecRec func(core.DecRec) (DecRec, error)
	dataToTermRef      func(DecRec) (termsem.SemRef, error)
)
