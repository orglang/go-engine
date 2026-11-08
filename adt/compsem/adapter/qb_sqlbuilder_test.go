package adapter

import (
	compsem "orglang/go-engine/adt/compsem"
	"testing"
)

func TestUpdateRef(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.updateRef(compsem.SemRefDS{})
	t.Log(sql)
}
