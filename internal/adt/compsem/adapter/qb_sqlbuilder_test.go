package adapter

import (
	"testing"
)

func TestUpdateRef(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.updateRef(SemRef{})
	t.Log(sql)
}
