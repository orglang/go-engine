package termdec

import (
	"orglang/go-engine/adt/semtype"

	"github.com/orglang/go-sdk/adt/pooldec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/pool/typedef:Msg.*
var (
	MsgToDecSpec   func(pooldec.DecSpec) (DecSpec, error)
	MsgFromDecSpec func(DecSpec) pooldec.DecSpec
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/descvar:Data.*
var (
	// goverter:map . DescRef | DataToDescRef
	DataToDecRec func(decRecDS) (DecRec, error)
	// goverter:autoMap DescRef
	DataFromDecRec func(DecRec) (decRecDS, error)
	// goverter:ignore DescRN
	DataToDescRef func(decRecDS) (semtype.TypeRef, error)
)
