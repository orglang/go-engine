package termsem

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/termsem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	MsgFromRef  func(SemRef) termsem.SemRef
	MsgFromRefs func([]SemRef) []termsem.SemRef
	MsgToRef    func(termsem.SemRef) (SemRef, error)
	MsgToRefs   func([]termsem.SemRef) ([]SemRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	DataFromRef  func(SemRef) SemRefDS
	DataFromRefs func([]SemRef) ([]SemRefDS, error)
	DataToRef    func(SemRefDS) (SemRef, error)
	DataToRefs   func([]SemRefDS) ([]SemRef, error)
	DataToRefMap func(map[uniqsym.ADT]SemRefDS) (map[uniqsym.ADT]SemRef, error)
)
