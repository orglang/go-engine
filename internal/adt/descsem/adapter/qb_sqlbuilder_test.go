package adapter

import "testing"

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.insertRec(semRec{})
	t.Log(sql)
}
