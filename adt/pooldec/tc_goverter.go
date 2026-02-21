package pooldec

import (
	"github.com/orglang/go-sdk/adt/pooldec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/xactdef:Msg.*
var (
	MsgToDecSpec   func(pooldec.DecSpec) (DecSpec, error)
	MsgFromDecSpec func(DecSpec) pooldec.DecSpec
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
var (
	// goverter:map . DescRef
	DataToDecRec func(decRecDS) (DecRec, error)
	// goverter:autoMap DescRef
	DataFromDecRec func(DecRec) (decRecDS, error)
)
