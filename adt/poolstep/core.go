package poolstep

import (
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/poolexp"
)

type StepSpec struct {
	ImplRef implsem.SemRef
	PoolES  poolexp.ExpSpec
}
