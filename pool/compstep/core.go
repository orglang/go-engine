package compstep

import (
	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/pool/termexp"
)

type StepSpec struct {
	CompRef compsem.SemRef
	PoolExp termexp.ExpSpec
}
