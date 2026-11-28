package adapter

import (
	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/uniqsym"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
var (
	DataFromRef  func(execroot.RootRef) RootRef
	DataFromRefs func([]execroot.RootRef) ([]RootRef, error)
	DataToRef    func(RootRef) (execroot.RootRef, error)
	DataToRefs   func([]RootRef) ([]execroot.RootRef, error)
	DataToRefMap func(map[uniqsym.ADT]RootRef) (map[uniqsym.ADT]execroot.RootRef, error)
)
