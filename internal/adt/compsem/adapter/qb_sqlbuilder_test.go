package adapter

import (
	"testing"

	compsem "orglang/go-engine/internal/adt/compsem"
)

func TestUpdateRef(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.updateRef(compsem.SemRefDS{})
	t.Log(sql)
}
