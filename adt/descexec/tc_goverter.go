package descexec

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/descexec"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
var (
	MsgFromRef  func(ExecRef) descexec.ExecRef
	MsgFromRefs func([]ExecRef) []descexec.ExecRef
	MsgToRef    func(descexec.ExecRef) (ExecRef, error)
	MsgToRefs   func([]descexec.ExecRef) ([]ExecRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
var (
	DataFromRef  func(ExecRef) (ExecRefDS, error)
	DataFromRefs func([]ExecRef) ([]ExecRefDS, error)
	DataToRef    func(ExecRefDS) (ExecRef, error)
	DataToRefs   func([]ExecRefDS) ([]ExecRef, error)
	// goverter:autoMap Ref
	DataFromRec func(ExecRec) (execRecDS, error)
	// goverter:map . Ref
	DataToRec      func(execRecDS) (ExecRec, error)
	DataToDescRefs func(map[uniqsym.ADT]ExecRefDS) (map[uniqsym.ADT]ExecRef, error)
)
