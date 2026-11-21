package adapter

import (
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
var (
	DataFromRef  func(typesem.SemRef) SemRef
	DataFromRefs func([]typesem.SemRef) ([]SemRef, error)
	DataToRef    func(SemRef) (typesem.SemRef, error)
	DataToRefs   func([]SemRef) ([]typesem.SemRef, error)
	DataToRefMap func(map[uniqsym.ADT]SemRef) (map[uniqsym.ADT]typesem.SemRef, error)
)
