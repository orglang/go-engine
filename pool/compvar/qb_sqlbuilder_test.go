package compvar

import (
	"testing"

	"orglang/go-engine/adt/compvar"
)

func TestInsert(t *testing.T) {
	qb := newSQLBuilder()
	rec := compvar.SampleVarRec()
	sql, args := qb.insertRec(poolStructVars, compvar.DataFromVarRec(rec))
	t.Log(sql)
	t.Log(args)
}
