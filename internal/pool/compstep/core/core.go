package core

import (
	"orglang/go-engine/internal/adt/comproot"
	compexp "orglang/go-engine/internal/pool/compexec/core"
)

type StepSpec struct {
	CompRef comproot.SemRef
	PoolExp compexp.ExpSpec
}
