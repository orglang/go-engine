package poolexec

import (
	"fmt"
	"testing"

	"orglang/go-engine/adt/implsem"
)

func TestInsert(t *testing.T) {
	qb := newSQLBuilder()
	rec := execRecDS{ImplID: "id-1"}
	sql, args := qb.insertRec(rec)
	fmt.Println(sql)
	fmt.Println(args)
}

func TestSelectSnap(t *testing.T) {
	qb := newSQLBuilder()
	sem := implsem.SemRefDS{ImplID: "id-1"}
	sql, args := qb.selectSnap(sem)
	fmt.Println(sql)
	fmt.Println(args)
}
