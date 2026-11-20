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
	DataFromRef  func(typesem.SemRef) SemRefDS
	DataFromRefs func([]typesem.SemRef) ([]SemRefDS, error)
	DataToRef    func(SemRefDS) (typesem.SemRef, error)
	DataToRefs   func([]SemRefDS) ([]typesem.SemRef, error)
	DataToRefMap func(map[uniqsym.ADT]SemRefDS) (map[uniqsym.ADT]typesem.SemRef, error)
)
