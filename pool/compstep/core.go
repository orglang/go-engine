package compstep

import (
	"orglang/go-engine/adt/semcomp"
	"orglang/go-engine/pool/termexp"
)

type StepSpec struct {
	CompRef semcomp.CompRef
	PoolExp termexp.ExpSpec
}
