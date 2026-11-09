package adapter

import (
	"slices"
	"testing"

	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(typeexp.StateDS{})
	wantSQL := "INSERT INTO pool_type_exps (exp_ck, sup_exp_ck, kind, spec) " +
		"VALUES ($1, $2, $3, $4) ON CONFLICT (exp_ck) DO NOTHING"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{int64(0), int64(0), typeexp.ExpKind(0), typeexp.ExpSpecDS{}}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestSelectRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL := qb.SelectRecByCK()
	wantSQL := "WITH RECURSIVE exp_tree AS ((SELECT top.exp_ck, top.sup_exp_ck, top.kind, top.spec " +
		"FROM pool_type_exps top WHERE top.exp_ck = $1) UNION ALL (SELECT sub.exp_ck, sub.sup_exp_ck, " +
		"sub.kind, sub.spec FROM exp_tree sup, pool_type_exps sub WHERE sub.sup_exp_ck = sup.exp_ck)) " +
		"SELECT * FROM exp_tree"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
}
