package poolcomm

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(connRecDS{})
	fmt.Println(sql)
}

func TestUpdateRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.updateRec(connModDS{CommON: sql.Null[int64]{Valid: true}})
	fmt.Println(sql)
}

func TestSelectSnap(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.selectSnap(connQryDS{ChnlID: sql.Null[string]{Valid: true}})
	fmt.Println(sql)
}
