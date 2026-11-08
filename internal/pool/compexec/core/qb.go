package core

import (
	"orglang/go-engine/internal/adt/compsem"
)

type QueryBuilder interface {
	InsertRec(ExecRecDS) (string, []any)
	SelectRecByRef(compsem.SemRefDS) (string, []any)
	SelectSnapByQN(string) (string, []any)
}
