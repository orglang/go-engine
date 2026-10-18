package core

import (
	"orglang/go-engine/adt/compsem"
)

type QueryBuilder interface {
	InsertRec(ExecRecDS) (string, []any)
	SelectRecByRef(compsem.SemRefDS) (string, []any)
	SelectSnapByQN(string) (string, []any)
}
