package compexec

import (
	"fmt"
	"orglang/go-engine/adt/compsem"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(execRec{})
	fmt.Println(sql)
}

func TestSelectSnap(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.selectRecByRef(compsem.SemRefDS{})
	fmt.Println(sql)
}
