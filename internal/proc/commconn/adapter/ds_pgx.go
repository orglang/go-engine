package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/connsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	commconn "orglang/go-engine/internal/proc/commconn/core"
	commturnadapter "orglang/go-engine/internal/proc/commturn/adapter"
	commturn "orglang/go-engine/internal/proc/commturn/core"
)

type daoPgx struct {
	log *slog.Logger
}

func NewDaoPgx(log *slog.Logger) commconn.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{log: log.With(name)}
}

// for compilation purposes
func newRepo() commconn.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec commconn.ConnRec) error {
	ref := connsem.DataFromRef(rec.ConnRef)
	const query = `
		INSERT INTO proc_comp_conns (conn_id, conn_rn, offset_nr)
		VALUES ($1, $2, $3)
	`
	_, err := conn.Pgx.Exec(ctx, query, ref.ConnID, ref.ConnRN, int64(rec.OffsetNr))
	if err != nil {
		dao.log.Error("Connection insertion failed", slog.Any("ref", rec.ConnRef))
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection insertion succeeded", slog.Any("ref", rec.ConnRef))
	return nil
}

func (dao *daoPgx) ModifyRec(ctx context.Context, conn tx.Conn, mod commconn.ConnMod) error {
	ensureConnection := func(ref connsem.SemRef) error {
		dto := connsem.DataFromRef(ref)
		const query = `
			INSERT INTO proc_comp_conns (conn_id, conn_rn, offset_nr)
			VALUES ($1, $2, 0)
			ON CONFLICT (conn_id) DO NOTHING
		`
		_, err := conn.Pgx.Exec(ctx, query, dto.ConnID, dto.ConnRN)
		return err
	}
	if mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty() {
		if err := ensureConnection(mod.ConnRef); err != nil {
			return err
		}
		ref := connsem.DataFromRef(mod.ConnRef)
		const query = `
			UPDATE proc_comp_conns
			SET offset_nr = $1
			WHERE conn_id = $2 AND conn_rn = $3 AND offset_nr = $4
		`
		tag, err := conn.Pgx.Exec(ctx, query, int64(mod.OffsetNr.Get()), ref.ConnID, ref.ConnRN, int64(mod.PrevOffsetNr))
		if err != nil {
			dao.log.Error("Connection offset update failed", slog.Any("ref", mod.ConnRef))
			return err
		}
		if tag.RowsAffected() == 0 {
			return commconn.ErrSnapshotConflict(mod.ConnRef.ConnRN, mod.PrevOffsetNr)
		}
	}

	nextRNs := make(map[string]seqnum.ADT)
	for _, turn := range mod.Turns {
		ref := turn.GetConnRef()
		if err := ensureConnection(ref); err != nil {
			return err
		}
		key := ref.ConnID.String()
		expectedRN, ok := nextRNs[key]
		if !ok {
			expectedRN = ref.ConnRN
		}
		const updateRevision = `
			UPDATE proc_comp_conns
			SET conn_rn = conn_rn + 1
			WHERE conn_id = $1 AND conn_rn = $2
		`
		tag, err := conn.Pgx.Exec(ctx, updateRevision, key, int64(expectedRN))
		if err != nil {
			dao.log.Error("Connection revision update failed", slog.Any("ref", ref))
			return err
		}
		if tag.RowsAffected() == 0 {
			return commconn.ErrRevisionConflict(expectedRN)
		}
		dto, err := commturnadapter.DataFromStepRec(turn)
		if err != nil {
			return err
		}
		dto.ConnRN = int64(expectedRN)
		const insertAct = `
			INSERT INTO proc_conn_acts (conn_id, conn_rn, comp_id, chnl_id, kind, exp)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = conn.Pgx.Exec(ctx, insertAct, dto.ConnID, dto.ConnRN, dto.CompID, dto.ChnlID, dto.K, dto.Exp)
		if err != nil {
			dao.log.Error("Act insertion failed", slog.Any("act", turn))
			return err
		}
		nextRNs[key] = expectedRN.Next()
	}
	return nil
}

func (dao *daoPgx) GetSnapByQry(ctx context.Context, conn tx.Conn, qry commconn.ConnQry) (commconn.ConnSnap, error) {
	ref := connsem.DataFromRef(qry.ConnRef)
	query := `
		WITH turns AS (
			SELECT turn.*
			FROM proc_conn_acts turn
			WHERE turn.conn_id = $1
	`
	args := []any{ref.ConnID}
	if qry.ChnlID != nil && !qry.ChnlID.IsEmpty() {
		query += " AND turn.chnl_id = $2"
		args = append(args, identity.ConvertToString(qry.ChnlID.Get()))
	}
	query += `
		)
		SELECT
			conn.conn_id,
			conn.conn_rn,
			conn.offset_nr,
			(SELECT array_agg(row(t.*) ORDER BY t.conn_rn) FROM turns t WHERE t.conn_rn > conn.offset_nr) AS turns
		FROM proc_comp_conns conn
		WHERE conn.conn_id = $1
	`
	rows, err := conn.Pgx.Query(ctx, query, args...)
	if err != nil {
		dao.log.Error("Connection snapshot query failed", slog.Any("ref", qry.ConnRef))
		return commconn.ConnSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[commconn.ConnSnapDS])
	if errors.Is(err, pgx.ErrNoRows) {
		return commconn.ConnSnap{ConnRef: qry.ConnRef}, nil
	}
	if err != nil {
		dao.log.Error("Connection snapshot scan failed", slog.Any("ref", qry.ConnRef))
		return commconn.ConnSnap{}, err
	}
	connID, err := identity.ConvertFromString(dto.ConnID)
	if err != nil {
		return commconn.ConnSnap{}, err
	}
	snap := commconn.ConnSnap{
		ConnRef:  connsem.SemRef{ConnID: connID, ConnRN: seqnum.ADT(dto.ConnRN)},
		OffsetNr: seqnum.ADT(dto.OffsetNr),
		Turns:    make([]commturn.TurnRec, 0, len(dto.Turns)),
	}
	for _, turnDTO := range dto.Turns {
		turn, err := commturnadapter.DataToStepRec(turnDTO)
		if err != nil {
			return commconn.ConnSnap{}, err
		}
		snap.Turns = append(snap.Turns, turn)
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot loaded", slog.Any("ref", qry.ConnRef))
	return snap, nil
}
