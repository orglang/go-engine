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

	conncore "orglang/go-engine/internal/proc/conn/core"
	connactadapter "orglang/go-engine/internal/proc/connact/adapter"
	connact "orglang/go-engine/internal/proc/connact/core"
)

type daoPgx struct {
	log *slog.Logger
}

func NewDaoPgx(log *slog.Logger) conncore.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{log: log.With(name)}
}

// for compilation purposes
func newRepo() conncore.Repo { return new(daoPgx) }

func (dao *daoPgx) ModifyRec(ctx context.Context, conn tx.Conn, mod conncore.ConnMod) error {
	var offsetNr any
	if mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty() {
		offsetNr = int64(mod.OffsetNr.Get())
	}

	var actNr any
	if len(mod.Acts) > 0 {
		// The provider CAS advances comp_rn exactly once for this operation.
		// Mirror that revision as the Act number, as on the Pool side.
		actNr = int64(mod.ProviderCompRN) + 1
	}

	if offsetNr != nil || actNr != nil {
		const updateConn = `
			UPDATE proc_comp_conns
			SET offset_nr = COALESCE($1, offset_nr),
			    act_nr = COALESCE($2, act_nr)
			WHERE conn_id = $3
		`
		tag, err := conn.Pgx.Exec(ctx, updateConn, offsetNr, actNr, identity.ConvertToString(mod.ConnID))
		if err != nil {
			dao.log.Error("Connection update failed", slog.Any("conn_id", mod.ConnID))
			return err
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("process connection missing: %v", mod.ConnID)
		}
	}

	for _, act := range mod.Acts {
		dto, err := connactadapter.DataFromStepRec(act)
		if err != nil {
			return err
		}
		dto.ActNr = int64(mod.ProviderCompRN) + 1
		const insertAct = `
			INSERT INTO proc_conn_acts (conn_id, act_nr, comp_id, comp_rn, chnl_id, kind, exp)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = conn.Pgx.Exec(ctx, insertAct, dto.ConnID, dto.ActNr, dto.CompID, dto.CompRN, dto.ChnlID, dto.K, dto.Exp)
		if err != nil {
			dao.log.Error("Act insertion failed", slog.Any("act", act))
			return err
		}
	}
	return nil
}

func (dao *daoPgx) GetSnapByQry(ctx context.Context, conn tx.Conn, qry conncore.ConnQry) (conncore.ConnSnap, error) {
	connID := identity.ConvertToString(qry.ConnID)
	query := `
		WITH acts AS (
			SELECT act.*
			FROM proc_conn_acts act
			WHERE act.conn_id = $1
	`
	args := []any{connID}
	if qry.ChnlID != nil && !qry.ChnlID.IsEmpty() {
		query += " AND act.chnl_id = $2"
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
			(SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM acts t WHERE t.act_nr > conn.offset_nr) AS acts
		FROM proc_comp_conns conn
		JOIN proc_comp_execs exec ON exec.comp_id = conn.provider_comp_id
		WHERE conn.conn_id = $1
	`
	rows, err := conn.Pgx.Query(ctx, query, args...)
	if err != nil {
		dao.log.Error("Connection snapshot query failed", slog.Any("conn_id", qry.ConnID))
		return conncore.ConnSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[conncore.ConnSnapDS])
	if errors.Is(err, pgx.ErrNoRows) {
		return conncore.ConnSnap{ConnID: qry.ConnID}, nil
	}
	if err != nil {
		dao.log.Error("Connection snapshot scan failed", slog.Any("conn_id", qry.ConnID))
		return conncore.ConnSnap{}, err
	}
	parsedConnID, err := identity.ConvertFromString(dto.ConnID)
	if err != nil {
		return conncore.ConnSnap{}, err
	}
	providerCompID, err := identity.ConvertFromString(dto.ProviderCompID)
	if err != nil {
		return conncore.ConnSnap{}, err
	}
	snap := conncore.ConnSnap{ConnID: parsedConnID, ProviderCompID: providerCompID, ProviderCompRN: seqnum.ADT(dto.ProviderCompRN), OffsetNr: seqnum.ADT(dto.OffsetNr), Acts: make([]connact.ActRec, 0, len(dto.Acts))}
	for _, actDTO := range dto.Acts {
		act, err := connactadapter.DataToStepRec(actDTO)
		if err != nil {
			return conncore.ConnSnap{}, err
		}
		snap.Acts = append(snap.Acts, act)
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot loaded", slog.Any("conn_id", qry.ConnID))
	return snap, nil
}
