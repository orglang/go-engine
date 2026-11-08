package adapter

import (
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/connsem"

	commturn "orglang/go-engine/internal/pool/commturn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/pool/termexp/adapter:Data.*
var (
	DataToConnRef func(commturn.TurnRecDS) (connsem.SemRef, error)
	// goverter:ignore CompRN
	DataToCompRef func(commturn.TurnRecDS) (compsem.SemRef, error)
	// goverter:map . ConnRef
	// goverter:map . CompRef
	// goverter:map Exp ValExp
	DataToPubRec func(commturn.TurnRecDS) (commturn.PubRec, error)
	// goverter:map . ConnRef
	// goverter:map . CompRef
	// goverter:map Exp ContExp
	DataToSubRec func(commturn.TurnRecDS) (commturn.SubRec, error)
)
