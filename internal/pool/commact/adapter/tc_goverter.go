package adapter

import (
	"orglang/go-engine/internal/adt/compsem"

	commact "orglang/go-engine/internal/pool/commact/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/pool/termexp/adapter:Data.*
var (
	DataToCompRef func(commact.ActRecDS) (compsem.SemRef, error)
	// goverter:map . CompRef
	// goverter:map Exp ValExp
	DataToPubRec func(commact.ActRecDS) (commact.PubRec, error)
	// goverter:map . CompRef
	// goverter:map Exp ContExp
	DataToSubRec func(commact.ActRecDS) (commact.SubRec, error)
)
