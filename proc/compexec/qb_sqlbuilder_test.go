package compexec

import "testing"

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(execRecDS{})
	t.Log(sql)
}
