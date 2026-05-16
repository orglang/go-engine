package commturn

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/compsem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/pool/termexp:Data.*
var (
	DataToCommRef func(TurnRecDS) (commsem.SemRef, error)
	// goverter:ignore CompRN
	DataToCompRef func(TurnRecDS) (compsem.SemRef, error)
	// goverter:map . CommRef
	// goverter:map . CompRef
	// goverter:map Exp ValExp
	DataToPubRec func(TurnRecDS) (PubRec, error)
	// goverter:map . CommRef
	// goverter:map . CompRef
	// goverter:map Exp ContExp
	DataToSubRec func(TurnRecDS) (SubRec, error)
)
