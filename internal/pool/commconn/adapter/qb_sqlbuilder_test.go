package adapter

import (
	"database/sql"
	"slices"
	"testing"

	commconn "orglang/go-engine/internal/pool/commconn/core"
	commturn "orglang/go-engine/internal/pool/commturn/core"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(commconn.ConnRecDS{})
	wantSQL := "INSERT INTO pool_comp_conns  (conn_id, provider_comp_id, provider_comp_rn, client_comp_id, client_comp_rn, offset_nr, act_nr) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", "", int64(0), "", int64(0), int64(0), int64(0)}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestInsertTurnRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertTurnRec(commturn.TurnRecDS{})
	wantSQL := "INSERT INTO pool_conn_acts  (conn_id, act_nr, comp_id, comp_rn, chnl_id, kind, exp) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	if len(gotArgs) != 7 {
		t.Errorf("got %d args, want 7: %v", len(gotArgs), gotArgs)
	}
}

func TestUpdateRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.UpdateRec(commconn.ConnModDS{
		ConnID:   "",
		OffsetNr: sql.Null[int64]{V: 3, Valid: true},
	})
	wantSQL := "UPDATE pool_comp_conns  SET offset_nr = $1 WHERE conn_id = $2"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{int64(3), ""}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestSelectSnap(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.SelectSnap(commconn.ConnQryDS{ChnlID: sql.Null[string]{Valid: true}})
	wantSQL := "WITH turns AS (SELECT turn.conn_id, turn.act_nr, turn.comp_id, turn.comp_rn, turn.chnl_id, turn.kind, turn.exp " +
		"FROM pool_conn_acts turn WHERE conn_id = $1 AND chnl_id = $2) " +
		"SELECT conn.conn_id, conn.provider_comp_id, exec.comp_rn AS provider_comp_rn, conn.offset_nr, conn.act_nr, " +
		"(SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM turns t WHERE t.act_nr > conn.offset_nr) AS turns " +
		"FROM pool_comp_conns conn JOIN pool_comp_execs exec ON exec.comp_id = conn.provider_comp_id WHERE conn.conn_id = $3"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", "", ""}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
