package adapter

import (
	"slices"
	"testing"

	pooltypeexp "orglang/go-engine/pool/typeexp/core"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(pooltypeexp.StateDS{})
	wantSQL := "INSERT INTO pool_type_exps (exp_vk, sup_exp_vk, kind, spec) " +
		"VALUES ($1, $2, $3, $4) ON CONFLICT (exp_vk) DO NOTHING"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{int64(0), int64(0), pooltypeexp.ExpKind(0), pooltypeexp.ExpSpecDS{}}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestSelectRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL := qb.SelectRecByVK()
	wantSQL := "WITH RECURSIVE exp_tree AS ((SELECT top.exp_vk, top.sup_exp_vk, top.kind, top.spec " +
		"FROM pool_type_exps top WHERE top.exp_vk = $1) UNION ALL (SELECT sub.exp_vk, sub.sup_exp_vk, " +
		"sub.kind, sub.spec FROM exp_tree sup, pool_type_exps sub WHERE sub.sup_exp_vk = sup.exp_vk)) " +
		"SELECT * FROM exp_tree"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
}
