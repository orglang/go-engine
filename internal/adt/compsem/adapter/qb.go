package adapter

import compsem "orglang/go-engine/internal/adt/compsem"

type queryBuilder interface {
	updateRef(compsem.SemRefDS) (string, []any)
}
