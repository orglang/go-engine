package core

import (
	"orglang/go-engine/internal/adt/comproot"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

type StepSpec struct {
	CompRef comproot.SemRef
	PoolExp termexp.ExpSpec
}
