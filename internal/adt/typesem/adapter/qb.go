package adapter

import typesem "orglang/go-engine/internal/adt/typesem"

type queryBuilder interface {
	updateRef(typesem.SemRefDS) (string, []any)
	selectRefByQN() string
}
