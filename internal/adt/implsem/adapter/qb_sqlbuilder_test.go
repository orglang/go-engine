package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/implsem"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.insertRec(implsem.SemRecDS{})
	t.Log(sql)
}
