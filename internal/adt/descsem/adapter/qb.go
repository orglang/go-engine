package adapter

import "orglang/go-engine/internal/adt/descsem"

type queryBuilder interface {
	insertRec(descsem.SemRecDS) (string, []any)
}
