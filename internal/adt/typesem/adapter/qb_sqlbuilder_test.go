package adapter

import "testing"

func TestUpdateRef(t *testing.T) {
	qb := newSQLBuilder("types", "descs")
	sql, _ := qb.updateRef(SemRef{})
	t.Log(sql)
}
