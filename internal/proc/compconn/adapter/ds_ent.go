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
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/compconn/core"
	connactadapter "orglang/go-engine/internal/proc/connact/adapter"
	connact "orglang/go-engine/internal/proc/connact/core"
	"orglang/go-engine/internal/proc/ent"
	entcompconn "orglang/go-engine/internal/proc/ent/compconn"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) core.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() core.Repo { return new(daoEnt) }

func (dao *daoEnt) ModifyRec(ctx context.Context, conn tx.Conn, mod core.ConnMod) error {
	refAttr := slog.Any("conn_id", mod.ConnID)
	var offsetNr *int64
	if mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty() {
		value := int64(mod.OffsetNr.Get())
		offsetNr = &value
	}
	var actNr *int64
	if len(mod.Acts) > 0 {
		value := int64(mod.ProviderCompRN) + 1
		actNr = &value
	}

	if offsetNr != nil || actNr != nil {
		update := dao.client.CompConn.Update().
			Where(entcompconn.IDEQ(identity.ConvertToString(mod.ConnID)))
		if offsetNr != nil {
			update.SetOffsetNr(*offsetNr)
		}
		if actNr != nil {
			update.SetActNr(*actNr)
		}
		affected, err := update.Save(entpgx.WithConn(ctx, conn.Pgx))
		if err != nil {
			dao.log.Error("Connection update failed", refAttr)
			return err
		}
		if affected == 0 {
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

func (dao *daoEnt) GetSnapByQry(ctx context.Context, conn tx.Conn, qry core.ConnQry) (core.ConnSnap, error) {
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
		return core.ConnSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ConnSnap])
	if errors.Is(err, pgx.ErrNoRows) {
		return core.ConnSnap{ConnID: qry.ConnID}, nil
	}
	if err != nil {
		dao.log.Error("Connection snapshot scan failed", slog.Any("conn_id", qry.ConnID))
		return core.ConnSnap{}, err
	}
	parsedConnID, err := identity.ConvertFromString(dto.ConnID)
	if err != nil {
		return core.ConnSnap{}, err
	}
	providerCompID, err := identity.ConvertFromString(dto.ProviderCompID)
	if err != nil {
		return core.ConnSnap{}, err
	}
	snap := core.ConnSnap{
		ConnID: parsedConnID, ProviderCompID: providerCompID,
		ProviderCompRN: seqnum.ADT(dto.ProviderCompRN), OffsetNr: seqnum.ADT(dto.OffsetNr),
		Acts: make([]connact.ActRec, 0, len(dto.Acts)),
	}
	for _, actDTO := range dto.Acts {
		act, err := connactadapter.DataToStepRec(actDTO)
		if err != nil {
			return core.ConnSnap{}, err
		}
		snap.Acts = append(snap.Acts, act)
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot loaded", slog.Any("conn_id", qry.ConnID))
	return snap, nil
}
