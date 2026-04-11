package poolexec

import (
	"fmt"
	"testing"

	"orglang/go-engine/adt/implsem"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.insertRec(execRecDS{})
	fmt.Println(sql)
}

func TestSelectSnap(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.selectSnap(implsem.SemRefDS{})
	fmt.Println(sql)
}

func TestSelectSnapByQN(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.selectSnapByQN("qn-1")
	fmt.Println(sql)
}
