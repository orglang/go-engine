package commexch

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(exchRecDS{})
	fmt.Println(sql)
}

func TestUpdateRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.updateRec(exchModDS{OffsetNr: sql.Null[int64]{Valid: true}})
	fmt.Println(sql)
}

func TestSelectSnap(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.selectSnap(exchQryDS{ChnlID: sql.Null[string]{Valid: true}})
	fmt.Println(sql)
}
