package compstep

import (
	"orglang/go-engine/adt/compsem"
	termexp "orglang/go-engine/pool/termexp/core"
)

type StepSpec struct {
	CompRef compsem.SemRef
	PoolExp termexp.ExpSpec
}
