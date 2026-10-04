package compsem

import "testing"

func TestUpdateRef(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.updateRef(SemRefDS{})
	t.Log(sql)
}
