package syndec

import (
	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
	"orglang/orglang/adt/revnum"
)

type DecRec struct {
	DecID identity.ADT
	DecRN revnum.ADT
	DecQN qualsym.ADT
}
