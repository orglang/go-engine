package adapter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

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
func newRepo() commconn.Repo { return new(daoPgx) }

func (dao *daoPgx) ModifyRec(ctx context.Context, conn tx.Conn, mod commconn.ConnMod) error {
	if mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty() {
		const updateOffset = `
			UPDATE proc_comp_conns
			SET offset_nr = $1
			WHERE conn_id = $2
		`
		tag, err := conn.Pgx.Exec(ctx, updateOffset, int64(mod.OffsetNr.Get()), identity.ConvertToString(mod.ConnID))
		if err != nil {
			dao.log.Error("Connection offset update failed", slog.Any("conn_id", mod.ConnID))
			return err
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("process connection missing: %v", mod.ConnID)
		}
	}
	for _, turn := range mod.Turns {
		connID := turn.GetConnID()
		switch turn.(type) {
		case commturn.PubRec, commturn.SubRec:
		default:
			panic(commturn.ErrRecTypeUnexpected(turn))
		}
		const nextAct = `
			UPDATE proc_comp_conns
			SET act_nr = act_nr + 1
			WHERE conn_id = $1
			RETURNING act_nr
		`
		var actNr int64
		err := conn.Pgx.QueryRow(ctx, nextAct, identity.ConvertToString(connID)).Scan(&actNr)
		if err != nil {
			return err
		}
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
			exec.comp_id AS provider_comp_id,
			exec.comp_rn AS provider_comp_rn,
			conn.offset_nr,
			conn.act_nr,
			(SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM turns t WHERE t.act_nr > conn.offset_nr) AS turns
		FROM proc_comp_conns conn
		JOIN proc_comp_execs exec ON exec.comp_id = conn.provider_comp_id
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
	providerCompID, err := identity.ConvertFromString(dto.ProviderCompID)
	if err != nil {
		return commconn.ConnSnap{}, err
	}
	snap := commconn.ConnSnap{ConnID: parsedConnID, ProviderCompID: providerCompID, ProviderCompRN: seqnum.ADT(dto.ProviderCompRN), OffsetNr: seqnum.ADT(dto.OffsetNr), Turns: make([]commturn.TurnRec, 0, len(dto.Turns))}
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
