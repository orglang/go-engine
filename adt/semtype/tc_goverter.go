package semtype

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/descsem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	MsgFromRef  func(TypeRef) descsem.SemRef
	MsgFromRefs func([]TypeRef) []descsem.SemRef
	MsgToRef    func(descsem.SemRef) (TypeRef, error)
	MsgToRefs   func([]descsem.SemRef) ([]TypeRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	DataFromRef  func(TypeRef) (SemRefDS, error)
	DataFromRefs func([]TypeRef) ([]SemRefDS, error)
	DataToRef    func(SemRefDS) (TypeRef, error)
	DataToRefs   func([]SemRefDS) ([]TypeRef, error)
	DataFromBind func(SemBind) (SemBindDS, error)
	DataToBind   func(SemBindDS) (SemBind, error)
	// goverter:autoMap DescRef
	DataFromRec func(SemRec) (semRecDS, error)
	// goverter:map . DescRef
	DataToRec    func(semRecDS) (SemRec, error)
	DataToRefMap func(map[uniqsym.ADT]SemRefDS) (map[uniqsym.ADT]TypeRef, error)
)
