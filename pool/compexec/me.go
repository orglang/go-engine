package compexec

import (
	compstep "orglang/go-engine/pool/compstep/core"
)

type Broker interface {
	Subscribe(api API)
	SendSpec(compstep.StepSpec) error
}
