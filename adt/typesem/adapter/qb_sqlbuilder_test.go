package adapter

import (
	typesem "orglang/go-engine/adt/typesem"
	"testing"
)

func TestUpdateRef(t *testing.T) {
	qb := newSQLBuilder("types", "descs")
	sql, _ := qb.updateRef(typesem.SemRefDS{})
	t.Log(sql)
}
