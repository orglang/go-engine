package compexec

import (
	"fmt"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(execRecDS{})
	fmt.Println(sql)
}
