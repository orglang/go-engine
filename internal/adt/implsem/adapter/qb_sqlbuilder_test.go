package adapter

import "testing"

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.insertRec(semRecDS{})
	t.Log(sql)
}
