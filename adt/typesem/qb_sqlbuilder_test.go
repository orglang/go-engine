package typesem

import (
	"testing"
)

func TestUpdateRef(t *testing.T) {
	qb := newSQLBuilder("types", "descs")
	sql, _ := qb.updateRef(SemRefDS{})
	t.Log(sql)
}
