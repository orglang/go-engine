package commturn

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/compsem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	DataToCommRef func(TurnRecDS) (commsem.SemRef, error)
	// goverter:ignore CompRN
	DataToCompRef func(TurnRecDS) (compsem.SemRef, error)
)
