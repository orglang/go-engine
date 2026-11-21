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
	DataFromRef  func(termsem.SemRef) SemRef
	DataFromRefs func([]termsem.SemRef) ([]SemRef, error)
	DataToRef    func(SemRef) (termsem.SemRef, error)
	DataToRefs   func([]SemRef) ([]termsem.SemRef, error)
	DataToRefMap func(map[uniqsym.ADT]SemRef) (map[uniqsym.ADT]termsem.SemRef, error)
)
