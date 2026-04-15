package compexec

import (
	"orglang/go-engine/pool/compstep"
)

type Exch interface {
	Subscribe(api API)
	SendSpec(compstep.StepSpec) error
}
