package poolimpl

import (
	"orglang/go-engine/adt/poolstep"
)

type Exch interface {
	Subscribe(api API)
	SendSpec(poolstep.StepSpec) error
}
