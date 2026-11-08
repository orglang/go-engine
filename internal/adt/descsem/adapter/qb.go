package adapter

import descsem "orglang/go-engine/adt/descsem"

type queryBuilder interface {
	insertRec(descsem.SemRecDS) (string, []any)
}
