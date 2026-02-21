package descsem

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/descsem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
var (
	MsgFromRef  func(SemRef) descsem.SemRef
	MsgFromRefs func([]SemRef) []descsem.SemRef
	MsgToRef    func(descsem.SemRef) (SemRef, error)
	MsgToRefs   func([]descsem.SemRef) ([]SemRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	DataFromRef  func(SemRef) (SemRefDS, error)
	DataFromRefs func([]SemRef) ([]SemRefDS, error)
	DataToRef    func(SemRefDS) (SemRef, error)
	DataToRefs   func([]SemRefDS) ([]SemRef, error)
	DataFromBind func(SemBind) (semBindDS, error)
	DataToBind   func(semBindDS) (SemBind, error)
	// goverter:autoMap Ref
	// goverter:map Bind.DescQN DescQN
	DataFromRec func(SemRec) (semRecDS, error)
	// goverter:map . Ref
	// goverter:map . Bind
	DataToRec    func(semRecDS) (SemRec, error)
	DataToRefMap func(map[uniqsym.ADT]SemRefDS) (map[uniqsym.ADT]SemRef, error)
)
