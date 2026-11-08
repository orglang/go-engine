package adapter

import (
	"testing"

	descsem "orglang/go-engine/internal/adt/descsem"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.insertRec(descsem.SemRecDS{})
	t.Log(sql)
}
