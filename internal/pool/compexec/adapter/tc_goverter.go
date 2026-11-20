package adapter

import (
	"orglang/go-engine/internal/adt/uniqsym"

	compexec "orglang/go-engine/internal/pool/compexec/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compvar/adapter:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compvar/adapter:Data.*
// goverter:extend DataToExecSnap1
var (
	// goverter:map . CompRef
	DataToExecRec func(ExecRecDS) (compexec.ExecRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(compexec.ExecRec) ExecRecDS
	// goverter:map . CompRef
	DataToExecSnap2 func(ExecSnap2DS) (compexec.ExecSnap2, error)
	// goverter:autoMap CompRef
	DataFromExecSnap2 func(compexec.ExecSnap2) ExecSnap2DS

	DataToRefMap  func(map[uniqsym.ADT]ExecSnap2DS) (map[uniqsym.ADT]compexec.ExecSnap2, error)
	DataToSnapMap func(map[uniqsym.ADT]ExecSnap1DS) (map[uniqsym.ADT]compexec.ExecSnap1, error)
)
