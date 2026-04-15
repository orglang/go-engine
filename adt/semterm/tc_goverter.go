package semterm

import (
	"orglang/go-engine/adt/uniqsym"

	"github.com/orglang/go-sdk/adt/semterm"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	MsgFromRef  func(TermRef) semterm.TermRef
	MsgFromRefs func([]TermRef) []semterm.TermRef
	MsgToRef    func(semterm.TermRef) (TermRef, error)
	MsgToRefs   func([]semterm.TermRef) ([]TermRef, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	DataFromRef  func(TermRef) TermRefDS
	DataFromRefs func([]TermRef) ([]TermRefDS, error)
	DataToRef    func(TermRefDS) (TermRef, error)
	DataToRefs   func([]TermRefDS) ([]TermRef, error)
	DataToRefMap func(map[uniqsym.ADT]TermRefDS) (map[uniqsym.ADT]TermRef, error)
)
