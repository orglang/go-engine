package compsem

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/compsem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	MsgFromRef  func(SemRef) compsem.SemRef
	MsgFromRefs func([]SemRef) []compsem.SemRef
	MsgToRef    func(compsem.SemRef) (SemRef, error)
	MsgToRefs   func([]compsem.SemRef) ([]SemRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	DataFromRef  func(SemRef) SemRefDS
	DataFromRefs func([]SemRef) ([]SemRefDS, error)
	DataToRef    func(SemRefDS) (SemRef, error)
	DataToRefs   func([]SemRefDS) ([]SemRef, error)
	DataToRefMap func(map[uniqsym.ADT]SemRefDS) (map[uniqsym.ADT]SemRef, error)
)
