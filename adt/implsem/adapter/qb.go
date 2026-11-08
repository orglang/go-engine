package adapter

import implsem "orglang/go-engine/adt/implsem"

type queryBuilder interface {
	insertRec(implsem.SemRecDS) (string, []any)
}
