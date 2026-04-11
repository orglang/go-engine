package implsem

import (
	"fmt"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(semRecDS{})
	fmt.Println(sql)
}

func TestUpdateRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.updateRec(SemRefDS{})
	fmt.Println(sql)
}
