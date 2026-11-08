package adapter

import (
	"testing"

	typesem "orglang/go-engine/internal/adt/typesem"
)

func TestUpdateRef(t *testing.T) {
	qb := newSQLBuilder("types", "descs")
	sql, _ := qb.updateRef(typesem.SemRefDS{})
	t.Log(sql)
}
