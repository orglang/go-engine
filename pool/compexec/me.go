package compexec

import (
	"orglang/go-engine/pool/compstep"
)

type Broker interface {
	Subscribe(api API)
	SendSpec(compstep.StepSpec) error
}
