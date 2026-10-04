package adapter

import (
	"slices"
	"testing"

	pooltypedef "orglang/go-engine/pool/typedef/core"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(pooltypedef.DefRecDS{
		TypeID: "type-id-1",
		TypeRN: 7,
		ExpVK:  42,
	})
	wantSQL := "INSERT INTO pool_type_defs  (type_id, type_rn, exp_vk) VALUES ($1, $2, $3)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"type-id-1", int64(7), int64(42)}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestSelectRecByQN(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	got := qb.SelectRecByQN()
	want := "SELECT def.type_id, def.type_rn, def.exp_vk FROM pool_type_defs def " +
		"JOIN pool_desc_binds bind ON bind.desc_id = def.type_id WHERE bind.desc_qn = $1"
	if got != want {
		t.Errorf("got sql %q, want %q", got, want)
	}
}
