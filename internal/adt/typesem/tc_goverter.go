package typesem

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/typesem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	MsgFromRef  func(SemRef) typesem.SemRef
	MsgFromRefs func([]SemRef) []typesem.SemRef
	MsgToRef    func(typesem.SemRef) (SemRef, error)
	MsgToRefs   func([]typesem.SemRef) ([]SemRef, error)
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
