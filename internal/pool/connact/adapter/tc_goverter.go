package adapter

import (
	"orglang/go-engine/internal/adt/comproot"

	connact "orglang/go-engine/internal/pool/connact/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/pool/termexp/adapter:Data.*
var (
	DataToCompRef func(ActRec) (comproot.SemRef, error)
	// goverter:map . CompRef
	// goverter:map Exp ValExp
	DataToPubRec func(ActRec) (connact.PubRec, error)
	// goverter:map . CompRef
	// goverter:map Exp ContExp
	DataToSubRec func(ActRec) (connact.SubRec, error)
)
