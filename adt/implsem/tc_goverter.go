package implsem

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/implsem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
var (
	MsgFromRef  func(SemRef) implsem.SemRef
	MsgFromRefs func([]SemRef) []implsem.SemRef
	MsgToRef    func(implsem.SemRef) (SemRef, error)
	MsgToRefs   func([]implsem.SemRef) ([]SemRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
var (
	DataFromRef  func(SemRef) (SemRefDS, error)
	DataFromRefs func([]SemRef) ([]SemRefDS, error)
	DataToRef    func(SemRefDS) (SemRef, error)
	DataToRefs   func([]SemRefDS) ([]SemRef, error)
	// goverter:autoMap Ref
	DataFromRec func(SemRec) (semRecDS, error)
	// goverter:map . Ref
	DataToRec    func(semRecDS) (SemRec, error)
	DataToRefMap func(map[uniqsym.ADT]SemRefDS) (map[uniqsym.ADT]SemRef, error)
)
