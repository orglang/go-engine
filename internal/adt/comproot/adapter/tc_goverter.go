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
	DataFromRef  func(comproot.SemRef) SemRef
	DataFromRefs func([]comproot.SemRef) ([]SemRef, error)
	DataToRef    func(SemRef) (comproot.SemRef, error)
	DataToRefs   func([]SemRef) ([]comproot.SemRef, error)
	DataToRefMap func(map[uniqsym.ADT]SemRef) (map[uniqsym.ADT]comproot.SemRef, error)
)
