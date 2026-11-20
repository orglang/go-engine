package adapter

import (
	"orglang/go-engine/internal/adt/typesem"
)


type queryBuilder interface {
	updateRef(SemRefDS) (string, []any)
	selectRefByQN() string
}
