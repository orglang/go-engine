package adapter

import compsem "orglang/go-engine/adt/compsem"

type queryBuilder interface {
	updateRef(compsem.SemRefDS) (string, []any)
}
