package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	commconn "orglang/go-engine/internal/proc/commconn/core"
	commturnadapter "orglang/go-engine/internal/proc/commturn/adapter"
	commturn "orglang/go-engine/internal/proc/commturn/core"
)

type actSequenceDS struct {
	ActNr int64 `db:"act_nr"`
}

type daoPgx struct {
	log *slog.Logger
}

func NewDaoPgx(log *slog.Logger) commconn.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{log: log.With(name)}
}

// for compilation purposes
func newRepo() commconn.Repo { return new(daoPgx) }

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec commconn.ConnRec) error {
	const query = `
		INSERT INTO proc_comp_conns (conn_id, provider_comp_id, provider_comp_rn, client_comp_id, client_comp_rn, offset_nr, act_nr)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := conn.Pgx.Exec(ctx, query,
		identity.ConvertToString(rec.ConnID),
		rec.ProviderCompRef.CompID.String(), int64(rec.ProviderCompRef.CompRN),
		rec.ClientCompRef.CompID.String(), int64(rec.ClientCompRef.CompRN),
		int64(rec.OffsetNr), int64(rec.ActNr),
	)
	if err != nil {
		dao.log.Error("Connection insertion failed", slog.Any("conn_id", rec.ConnID))
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection insertion succeeded", slog.Any("conn_id", rec.ConnID))
	return nil
}

func (dao *daoPgx) ModifyRec(ctx context.Context, conn tx.Conn, mod commconn.ConnMod) error {
	if mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty() {
		const updateOffset = `
			UPDATE proc_comp_conns
			SET offset_nr = $1
			WHERE conn_id = $2 AND offset_nr = $3
		`
		tag, err := conn.Pgx.Exec(ctx, updateOffset, int64(mod.OffsetNr.Get()), identity.ConvertToString(mod.ConnID), int64(mod.PrevOffsetNr))
		if err != nil {
			dao.log.Error("Connection offset update failed", slog.Any("conn_id", mod.ConnID))
			return err
		}
		if tag.RowsAffected() == 0 {
			return commconn.ErrSnapshotConflict(mod.PrevOffsetNr)
		}
	}
	for _, turn := range mod.Turns {
		connID := turn.GetConnID()
		compRef := turnCompRef(turn)
		// Materialize both endpoints from the authoritative pool Connection.
		const ensureConnection = `
			INSERT INTO proc_comp_conns (
				conn_id, provider_comp_id, provider_comp_rn,
				client_comp_id, client_comp_rn, offset_nr, act_nr
			)
			SELECT
				conn_id, provider_comp_id, provider_comp_rn,
				client_comp_id, client_comp_rn, 0, 0
			FROM pool_comp_conns
			WHERE conn_id = $1
			ON CONFLICT (conn_id) DO UPDATE SET
				provider_comp_id = EXCLUDED.provider_comp_id,
				provider_comp_rn = EXCLUDED.provider_comp_rn,
				client_comp_id = EXCLUDED.client_comp_id,
				client_comp_rn = EXCLUDED.client_comp_rn
			WHERE proc_comp_conns.provider_comp_id = EXCLUDED.provider_comp_id
				AND proc_comp_conns.provider_comp_rn = EXCLUDED.provider_comp_rn
				AND proc_comp_conns.client_comp_id = EXCLUDED.client_comp_id
				AND proc_comp_conns.client_comp_rn = EXCLUDED.client_comp_rn
		`
		if _, err := conn.Pgx.Exec(ctx, ensureConnection, identity.ConvertToString(connID)); err != nil {
			return err
		}
		var nextAct string
		switch turn.(type) {
		case commturn.PubRec:
			nextAct = `
				UPDATE proc_comp_conns
				SET act_nr = act_nr + 1
				WHERE conn_id = $1 AND provider_comp_id = $2 AND provider_comp_rn = $3
				RETURNING act_nr
			`
		case commturn.SubRec:
			nextAct = `
				UPDATE proc_comp_conns
				SET act_nr = act_nr + 1
				WHERE conn_id = $1 AND client_comp_id = $2 AND client_comp_rn = $3
				RETURNING act_nr
			`
		default:
			panic(commturn.ErrRecTypeUnexpected(turn))
		}
		rows, err := conn.Pgx.Query(ctx, nextAct, identity.ConvertToString(connID), compRef.CompID.String(), int64(compRef.CompRN))
		if err != nil {
			return err
		}
		act, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[actSequenceDS])
		rows.Close()
		if err != nil {
			return err
		}
		actNr := act.ActNr
		dto, err := commturnadapter.DataFromStepRec(turn)
		if err != nil {
			return err
		}
		dto.ActNr = actNr
		const insertAct = `
			INSERT INTO proc_conn_acts (conn_id, act_nr, comp_id, comp_rn, chnl_id, kind, exp)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = conn.Pgx.Exec(ctx, insertAct, dto.ConnID, dto.ActNr, dto.CompID, dto.CompRN, dto.ChnlID, dto.K, dto.Exp)
		if err != nil {
			dao.log.Error("Act insertion failed", slog.Any("act", turn))
			return err
		}
	}
	return nil
}

func turnCompRef(turn commturn.TurnRec) compsem.SemRef {
	switch rec := turn.(type) {
	case commturn.PubRec:
		return rec.CompRef
	case commturn.SubRec:
		return rec.CompRef
	default:
		panic(commturn.ErrRecTypeUnexpected(turn))
	}
}

func (dao *daoPgx) GetSnapByQry(ctx context.Context, conn tx.Conn, qry commconn.ConnQry) (commconn.ConnSnap, error) {
	connID := identity.ConvertToString(qry.ConnID)
	query := `
		WITH turns AS (
			SELECT turn.*
			FROM proc_conn_acts turn
			WHERE turn.conn_id = $1
	`
	args := []any{connID}
	if qry.ChnlID != nil && !qry.ChnlID.IsEmpty() {
		query += " AND turn.chnl_id = $2"
		args = append(args, identity.ConvertToString(qry.ChnlID.Get()))
	}
	query += `
		)
		SELECT
			conn.conn_id,
			conn.offset_nr,
			conn.act_nr,
			(SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM turns t WHERE t.act_nr > conn.offset_nr) AS turns
		FROM proc_comp_conns conn
		WHERE conn.conn_id = $1
	`
	rows, err := conn.Pgx.Query(ctx, query, args...)
	if err != nil {
		dao.log.Error("Connection snapshot query failed", slog.Any("conn_id", qry.ConnID))
		return commconn.ConnSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[commconn.ConnSnapDS])
	if errors.Is(err, pgx.ErrNoRows) {
		return commconn.ConnSnap{ConnID: qry.ConnID}, nil
	}
	if err != nil {
		dao.log.Error("Connection snapshot scan failed", slog.Any("conn_id", qry.ConnID))
		return commconn.ConnSnap{}, err
	}
	parsedConnID, err := identity.ConvertFromString(dto.ConnID)
	if err != nil {
		return commconn.ConnSnap{}, err
	}
	snap := commconn.ConnSnap{ConnID: parsedConnID, OffsetNr: seqnum.ADT(dto.OffsetNr), Turns: make([]commturn.TurnRec, 0, len(dto.Turns))}
	for _, turnDTO := range dto.Turns {
		turn, err := commturnadapter.DataToStepRec(turnDTO)
		if err != nil {
			return commconn.ConnSnap{}, err
		}
		snap.Turns = append(snap.Turns, turn)
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot loaded", slog.Any("conn_id", qry.ConnID))
	return snap, nil
}
