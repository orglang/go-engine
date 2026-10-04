package termdec

import "testing"

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(decRecDS{})
	t.Log(sql)
}
