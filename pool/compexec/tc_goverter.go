package compexec

import (
	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/uniqsym"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/compvar:Convert.*
// goverter:extend orglang/go-engine/adt/compvar:Data.*
// goverter:extend DataToExecLiabSnap
var (
	// goverter:map . CompRef | dataToSemRef
	// goverter:ignore TermQN
	DataToExecRec func(execRecDS) (ExecRec, error)
	// goverter:autoMap CompRef
	DataFromExecRec func(ExecRec) execRecDS

	DataToRefMap  func(map[uniqsym.ADT]execRecDS) (map[uniqsym.ADT]ExecRec, error)
	DataToSnapMap func(map[uniqsym.ADT]execSnapDS) (map[uniqsym.ADT]ExecSnap1, error)

	dataToSemRef func(execRecDS) (compsem.SemRef, error)
)
