package commexch

import (
	"database/sql"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	query, _ := qb.insertRec(exchRecDS{})
	t.Log(query)
}

func TestUpdateRec(t *testing.T) {
	qb := newSQLBuilder()
	query, _ := qb.updateRec(exchModDS{OffsetNr: sql.Null[int64]{Valid: true}})
	t.Log(query)
}

func TestSelectSnap(t *testing.T) {
	qb := newSQLBuilder()
	query, _ := qb.selectSnap(exchQryDS{ChnlID: sql.Null[string]{Valid: true}})
	t.Log(query)
}
