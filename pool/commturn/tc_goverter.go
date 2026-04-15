package commturn

import (
	"orglang/go-engine/adt/semcomm"
	"orglang/go-engine/adt/semcomp"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	DataToCommRef func(TurnRecDS) (semcomm.CommRef, error)
	// goverter:ignore ImplRN
	DataToCompRef func(TurnRecDS) (semcomp.CompRef, error)
)
