package compexec

import (
	"fmt"
	"orglang/go-engine/adt/compsem"
	"testing"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(execRec{})
	t.Log(sql)
}

func TestSelectSnap(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.selectRecByRef(compsem.SemRefDS{})
	t.Log(sql)
}
