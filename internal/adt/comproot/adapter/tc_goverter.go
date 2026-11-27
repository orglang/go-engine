package adapter

import (
	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/uniqsym"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
var (
	DataFromRef  func(comproot.RootRef) RootRef
	DataFromRefs func([]comproot.RootRef) ([]RootRef, error)
	DataToRef    func(RootRef) (comproot.RootRef, error)
	DataToRefs   func([]RootRef) ([]comproot.RootRef, error)
	DataToRefMap func(map[uniqsym.ADT]RootRef) (map[uniqsym.ADT]comproot.RootRef, error)
)
