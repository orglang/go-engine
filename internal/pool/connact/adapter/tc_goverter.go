package adapter

import (
	"orglang/go-engine/internal/adt/compsem"

	connact "orglang/go-engine/internal/pool/connact/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/pool/termexp/adapter:Data.*
var (
	DataToCompRef func(connact.ActRecDS) (compsem.SemRef, error)
	// goverter:map . CompRef
	// goverter:map Exp ValExp
	DataToPubRec func(connact.ActRecDS) (connact.PubRec, error)
	// goverter:map . CompRef
	// goverter:map Exp ContExp
	DataToSubRec func(connact.ActRecDS) (connact.SubRec, error)
)
