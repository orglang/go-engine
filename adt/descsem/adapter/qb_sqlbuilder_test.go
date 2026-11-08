package adapter

import (
	descsem "orglang/go-engine/adt/descsem"
	"fmt"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder("foo")
	sql, _ := qb.insertRec(descsem.SemRecDS{})
	fmt.Println(sql)
}
