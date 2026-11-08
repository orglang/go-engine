package adapter

import (
	"testing"

	implsem "orglang/go-engine/adt/implsem"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.insertRec(implsem.SemRecDS{})
	t.Log(sql)
}
