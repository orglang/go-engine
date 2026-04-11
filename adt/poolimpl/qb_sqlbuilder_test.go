package poolimpl

import (
	"fmt"
	"testing"

	"orglang/go-engine/adt/implsem"
)

func TestSelectSnap(t *testing.T) {
	qb := newSQLBuilder()
	sql, _ := qb.selectRecByRef(implsem.SemRefDS{})
	fmt.Println(sql)
}
