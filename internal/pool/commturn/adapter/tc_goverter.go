package adapter

import (
	"orglang/go-engine/internal/adt/compsem"

	commturn "orglang/go-engine/internal/pool/commturn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/pool/termexp/adapter:Data.*
var (
	DataToCompRef func(commturn.ActRecDS) (compsem.SemRef, error)
	// goverter:map . CompRef
	// goverter:map Exp ValExp
	DataToPubRec func(commturn.ActRecDS) (commturn.PubRec, error)
	// goverter:map . CompRef
	// goverter:map Exp ContExp
	DataToSubRec func(commturn.ActRecDS) (commturn.SubRec, error)
)
