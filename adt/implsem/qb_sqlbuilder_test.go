package implsem

import "testing"

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.insertRec(SemRecDS{})
	t.Log(sql)
}
