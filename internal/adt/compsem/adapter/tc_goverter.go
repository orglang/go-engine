package adapter

import (
	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/uniqsym"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
var (
	DataFromRef  func(compsem.SemRef) SemRef
	DataFromRefs func([]compsem.SemRef) ([]SemRef, error)
	DataToRef    func(SemRef) (compsem.SemRef, error)
	DataToRefs   func([]SemRef) ([]compsem.SemRef, error)
	DataToRefMap func(map[uniqsym.ADT]SemRef) (map[uniqsym.ADT]compsem.SemRef, error)
)
