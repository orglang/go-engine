package adapter

import (
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/uniqsym"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
var (
	DataFromRef  func(termsem.SemRef) SemRefDS
	DataFromRefs func([]termsem.SemRef) ([]SemRefDS, error)
	DataToRef    func(SemRefDS) (termsem.SemRef, error)
	DataToRefs   func([]SemRefDS) ([]termsem.SemRef, error)
	DataToRefMap func(map[uniqsym.ADT]SemRefDS) (map[uniqsym.ADT]termsem.SemRef, error)
)
