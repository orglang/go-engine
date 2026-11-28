package adapter

import (
	"orglang/go-engine/internal/adt/uniqsym"

	execeval "orglang/go-engine/internal/pool/execeval/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/adt/execvar/adapter:Convert.*
// goverter:extend orglang/go-engine/internal/adt/execvar/adapter:Data.*
// goverter:extend DataToEvalSnap1
var (
	// goverter:map . RootRef
	DataToEvalRec func(EvalRec) (execeval.EvalRec, error)
	// goverter:autoMap RootRef
	DataFromEvalRec func(execeval.EvalRec) EvalRec
	// goverter:map . RootRef
	DataToEvalSnap2 func(EvalSnap2) (execeval.EvalSnap2, error)
	// goverter:autoMap RootRef
	DataFromEvalSnap2 func(execeval.EvalSnap2) EvalSnap2

	DataToRefMap  func(map[uniqsym.ADT]EvalSnap2) (map[uniqsym.ADT]execeval.EvalSnap2, error)
	DataToSnapMap func(map[uniqsym.ADT]EvalSnap1) (map[uniqsym.ADT]execeval.EvalSnap1, error)
)
