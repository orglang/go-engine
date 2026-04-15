package compexec

import (
	"fmt"
	"orglang/go-engine/adt/semcomp"
	"testing"
)

func TestSelectSnap(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.selectRecByRef(semcomp.CompRefDS{})
	fmt.Println(sql)
}
