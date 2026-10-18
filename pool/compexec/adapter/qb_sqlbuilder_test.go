package adapter

import (
	"slices"
	"testing"

	"orglang/go-engine/adt/compsem"

	compexec "orglang/go-engine/pool/compexec/core"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(compexec.ExecRecDS{})
	wantSQL := "INSERT INTO pool_comp_execs  (comp_id, comp_rn, liab_mode) VALUES ($1, $2, $3)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", int64(0), int16(0)}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestSelectRecByRef(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.SelectRecByRef(compsem.SemRefDS{})
	wantSQL := "WITH struct_vars AS (SELECT var.comp_id, var.comp_rn, var.comm_id, var.chnl_id, " +
		"var.chnl_ph, var.exp_vk, var.side FROM pool_struct_vars var WHERE comp_id = $1), " +
		"linear_vars AS (SELECT var.comp_id, var.comp_rn, var.comm_id, var.chnl_id, " +
		"var.chnl_ph, var.exp_vk, var.side FROM pool_linear_vars var WHERE comp_id = $2) " +
		"SELECT exec.comp_id, exec.comp_rn, exec.liab_mode, " +
		"(SELECT array_agg(row(r.*)) FROM struct_vars r) AS struct_vars, " +
		"(SELECT array_agg(row(r.*)) FROM linear_vars r) AS linear_vars " +
		"FROM pool_comp_execs exec WHERE exec.comp_id = $3"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", "", ""}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestSelectSnapByQN(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.SelectSnapByQN("someQN")
	wantSQL := "SELECT sem.comp_id, sem.comp_rn, exec.liab_mode, " +
		"struct_var.comp_id, struct_var.comp_rn, struct_var.comm_id, struct_var.chnl_id, " +
		"struct_var.chnl_ph, struct_var.exp_vk, struct_var.side, " +
		"linear_var.comp_id, linear_var.comp_rn, linear_var.comm_id, linear_var.chnl_id, " +
		"linear_var.chnl_ph, linear_var.exp_vk, linear_var.side " +
		"FROM pool_comp_execs  JOIN pool_comp_execs  ON exec.comp_id = sem.comp_id " +
		"JOIN pool_struct_vars  ON struct_var.comp_id = sem.comp_id AND struct_var.side = 1 AND exec.mode = 1 " +
		"JOIN pool_linear_vars  ON linear_var.comp_id = sem.comp_id AND linear_var.side = 1 AND exec.mode = 2 " +
		"JOIN pool_impl_binds  ON bind.impl_id = sem.comp_id " +
		"WHERE bind.impl_qn = $1 ORDER BY struct_var.comp_rn DESC, linear_var.comp_rn DESC LIMIT $2"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"someQN", 1}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
