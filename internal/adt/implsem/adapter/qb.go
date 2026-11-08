package adapter

import "orglang/go-engine/internal/adt/implsem"

type queryBuilder interface {
	insertRec(implsem.SemRecDS) (string, []any)
}
