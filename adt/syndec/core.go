package syndec

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/uniqsym"
)

type DecRec struct {
	DecID identity.ADT
	DecRN revnum.ADT
	DecQN uniqsym.ADT
}
