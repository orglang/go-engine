package adapter

import (
	"testing"

	commturn "orglang/go-engine/pool/commturn/core"
)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(commturn.TurnRecDS{})
	wantSQL := "INSERT INTO pool_comm_turns  (comm_id, comm_rn, comp_id, chnl_id, kind, exp) VALUES ($1, $2, $3, $4, $5, $6)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	if len(gotArgs) != 6 {
		t.Errorf("got %d args, want 6: %v", len(gotArgs), gotArgs)
	}
}
