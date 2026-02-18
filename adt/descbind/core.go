package descbind

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/uniqsym"
)

type BindRec struct {
	// description qualified name
	DescQN uniqsym.ADT
	// description id
	DescID identity.ADT
}
