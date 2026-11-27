package adapter

import (
	"orglang/go-engine/internal/adt/comproot"
	execeval "orglang/go-engine/internal/proc/execeval/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compvar/adapter:Data.*
var (
	// goverter:map . CompRef | dataToSemRef
	DataToExecRec func(EvalRec) (execeval.EvalRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(execeval.EvalRec) EvalRec
	// goverter:ignore CompRN
	dataToSemRef func(EvalRec) (comproot.SemRef, error)
)
