package adapter

import (
	"orglang/go-engine/internal/adt/commsem"
	"orglang/go-engine/internal/adt/compsem"

	commturn "orglang/go-engine/internal/pool/commturn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/pool/termexp/adapter:Data.*
var (
	DataToCommRef func(commturn.TurnRecDS) (commsem.SemRef, error)
	// goverter:ignore CompRN
	DataToCompRef func(commturn.TurnRecDS) (compsem.SemRef, error)
	// goverter:map . CommRef
	// goverter:map . CompRef
	// goverter:map Exp ValExp
	DataToPubRec func(commturn.TurnRecDS) (commturn.PubRec, error)
	// goverter:map . CommRef
	// goverter:map . CompRef
	// goverter:map Exp ContExp
	DataToSubRec func(commturn.TurnRecDS) (commturn.SubRec, error)
)
