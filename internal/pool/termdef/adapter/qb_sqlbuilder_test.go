package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/adt/termvar"

	termdef "orglang/go-engine/pool/termdef/core"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(termdef.DefRecDS{})
	wantSQL := "INSERT INTO pool_term_decs  (term_id, term_rn, liab_var, asset_vars) VALUES ($1, $2, $3, $4)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", int64(0), termvar.VarRecDS{}, []termvar.VarRecDS(nil)}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestSelectRecByQN(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.SelectRecByQN("qn")
	wantSQL := "SELECT def.term_id, def.term_rn, def.liab_var, def.asset_vars FROM pool_term_decs def " +
		"JOIN pool_desc_binds bind ON bind.desc_id = def.term_id WHERE bind.desc_qn = $1"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"qn"}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
