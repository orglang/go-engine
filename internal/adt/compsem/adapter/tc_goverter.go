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
	DataFromRef  func(compsem.SemRef) SemRefDS
	DataFromRefs func([]compsem.SemRef) ([]SemRefDS, error)
	DataToRef    func(SemRefDS) (compsem.SemRef, error)
	DataToRefs   func([]SemRefDS) ([]compsem.SemRef, error)
	DataToRefMap func(map[uniqsym.ADT]SemRefDS) (map[uniqsym.ADT]compsem.SemRef, error)
)
