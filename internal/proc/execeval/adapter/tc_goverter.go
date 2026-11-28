package adapter

import (
	"orglang/go-engine/internal/adt/execroot"
	execeval "orglang/go-engine/internal/proc/execeval/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/execvar/adapter:Data.*
var (
	// goverter:map . RootRef | dataToRootRef
	DataToEvalRec func(EvalRec) (execeval.EvalRec, error)
	// goverter:autoMap RootRef
	DataFromEvalRec func(execeval.EvalRec) EvalRec
	// goverter:ignore ExecRN
	dataToRootRef func(EvalRec) (execroot.RootRef, error)
)
