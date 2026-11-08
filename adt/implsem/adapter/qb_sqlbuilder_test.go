package adapter

import (
	implsem "orglang/go-engine/adt/implsem"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.insertRec(implsem.SemRecDS{})
	t.Log(sql)
}
