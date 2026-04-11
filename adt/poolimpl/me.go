package poolimpl

import (
	"orglang/go-engine/adt/poolstep"
)

type Exch interface {
	SendSpec(poolstep.StepSpec) error
}
