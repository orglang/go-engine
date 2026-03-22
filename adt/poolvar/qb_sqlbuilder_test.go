package poolvar

import (
	"fmt"
	"testing"

	"orglang/go-engine/adt/implvar"
)

func TestInsert(t *testing.T) {
	qb := newSQLBuilder()
	rec := implvar.SampleVarRec()
	sql, args := qb.insertRec(poolStructVars, implvar.DataFromVarRec(rec))
	fmt.Println(sql)
	fmt.Println(args)
}
