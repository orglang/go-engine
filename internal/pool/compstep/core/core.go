package core

import (
	"orglang/go-engine/internal/adt/compsem"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

type StepSpec struct {
	CompRef compsem.SemRef
	PoolExp termexp.ExpSpec
}
