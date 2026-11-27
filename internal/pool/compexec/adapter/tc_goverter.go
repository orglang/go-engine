package adapter

import (
	"orglang/go-engine/internal/adt/uniqsym"

	compexec "orglang/go-engine/internal/pool/compexec/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compvar/adapter:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compvar/adapter:Data.*
// goverter:extend DataToExecSnap1
var (
	// goverter:map . CompRef
	DataToExecRec func(EvalRec) (compexec.EvalRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(compexec.EvalRec) EvalRec
	// goverter:map . CompRef
	DataToExecSnap2 func(EvalSnap2) (compexec.EvalSnap2, error)
	// goverter:autoMap CompRef
	DataFromExecSnap2 func(compexec.EvalSnap2) EvalSnap2

	DataToRefMap  func(map[uniqsym.ADT]EvalSnap2) (map[uniqsym.ADT]compexec.EvalSnap2, error)
	DataToSnapMap func(map[uniqsym.ADT]EvalSnap1) (map[uniqsym.ADT]compexec.EvalSnap1, error)
)
