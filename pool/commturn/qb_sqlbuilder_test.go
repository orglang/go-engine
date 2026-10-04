package commturn

import (
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(TurnRecDS{})
	t.Log(sql)
}
