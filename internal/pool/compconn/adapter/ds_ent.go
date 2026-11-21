package adapter

import (
	"context"
	"database/sql"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	entcompconn "orglang/go-engine/internal/pool/ent/compconn"

	compconn "orglang/go-engine/internal/pool/compconn/core"
	connactadapter "orglang/go-engine/internal/pool/connact/adapter"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) compconn.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() compconn.Repo { return new(daoEnt) }

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec compconn.ConnRec) error {
	dto := DataFromRec(rec)
	refAttr := slog.Any("conn_id", rec.ConnID)
	txCtx := entpgx.WithConn(ctx, conn.Pgx)
	_, err := dao.client.CompConn.Create().
		SetID(dto.ConnID).
		SetProviderCompID(dto.ProviderCompID).
		SetProviderCompRn(dto.ProviderCompRN).
		SetClientCompID(dto.ClientCompID).
		SetClientCompRn(dto.ClientCompRN).
		SetOffsetNr(dto.OffsetNr).
		SetActNr(dto.ActNr).
		Save(txCtx)
	if err != nil {
		dao.log.Error("Connection insertion failed", refAttr)
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection insertion succeeded", refAttr)
	return nil
}

func (dao *daoEnt) ModifyRec(ctx context.Context, conn tx.Conn, mod compconn.ConnMod) error {
	refAttr := slog.Any("conn_id", mod.ConnID)
	if (mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty()) || len(mod.Acts) > 0 {
		dto := DataFromMod(mod)
		if len(mod.Acts) > 0 {
			// The provider CAS advances comp_rn exactly once for this operation.
			// Mirror that revision as the Act number; Take appends at most one Act per operation.
			dto.ActNr = sql.Null[int64]{V: int64(mod.ProviderCompRN) + 1, Valid: true}
		}

		txCtx := entpgx.WithConn(ctx, conn.Pgx)
		update := dao.client.CompConn.Update().
			Where(entcompconn.IDEQ(dto.ConnID))
		if dto.OffsetNr.Valid {
			update.SetOffsetNr(dto.OffsetNr.V)
		}
		if dto.ActNr.Valid {
			update.SetActNr(dto.ActNr.V)
		}
		affected, err := update.Save(txCtx)
		if err != nil {
			dao.log.Error("Connection offset update failed", refAttr)
			return err
		}
		if affected == 0 {
			return compconn.ErrMissingConn(mod.ConnID)
		}
	}

	for _, act := range mod.Acts {
		dto := connactadapter.DataFromStepRec(act)
		dto.ActNr = int64(mod.ProviderCompRN) + 1
		_, err := conn.Pgx.Exec(ctx, insertActRec,
			dto.ConnID, dto.ActNr, dto.CompID, dto.CompRN, dto.ChnlID, dto.K, dto.Exp,
		)
		if err != nil {
			dao.log.Error("Act insertion failed", refAttr)
			return err
		}
	}

	dao.log.Log(ctx, lf.LevelTrace, "Connection modification succeeded", refAttr)
	return nil
}

func (dao *daoEnt) GetSnapByQry(ctx context.Context, conn tx.Conn, qry compconn.ConnQry) (compconn.ConnSnap, error) {
	refAttr := slog.Any("conn_id", qry.ConnID)
	qryDTO := DataFromQry(qry)
	querySQL := selectSnap
	args := []any{qryDTO.ConnID}
	if qryDTO.ChnlID.Valid {
		querySQL = selectSnapByChannel
		args = append(args, qryDTO.ChnlID.V)
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot selection started", slog.Any("qry", qryDTO))
	rows, err := conn.Pgx.Query(ctx, querySQL, args...)
	if err != nil {
		dao.log.Error("Connection snapshot query failed", refAttr)
		return compconn.ConnSnap{}, err
	}
	defer rows.Close()

	snapDTO, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ConnSnap])
	if err != nil {
		dao.log.Error("Connection snapshot scan failed", refAttr)
		return compconn.ConnSnap{}, err
	}
	snap, err := DataToSnap(snapDTO)
	if err != nil {
		dao.log.Error("Connection snapshot conversion failed", refAttr)
		return compconn.ConnSnap{}, err
	}
	providerCompID, err := identity.ConvertFromString(snapDTO.ProviderCompID)
	if err != nil {
		return compconn.ConnSnap{}, err
	}
	snap.ProviderCompID = providerCompID
	snap.ProviderCompRN = seqnum.ADT(snapDTO.ProviderCompRN)
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot loaded", slog.Any("dto", snapDTO))
	return snap, nil
}

const (
	insertActRec = `INSERT INTO pool_conn_acts
		(conn_id, act_nr, comp_id, comp_rn, chnl_id, kind, exp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	selectSnap = `WITH acts AS (
		SELECT act.conn_id, act.act_nr, act.comp_id, act.comp_rn, act.chnl_id, act.kind, act.exp
		FROM pool_conn_acts act WHERE act.conn_id = $1
	)
	SELECT conn.conn_id, conn.provider_comp_id, exec.comp_rn AS provider_comp_rn,
		conn.offset_nr, conn.act_nr,
		(SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM acts t WHERE t.act_nr > conn.offset_nr) AS acts
	FROM pool_comp_conns conn
	JOIN pool_comp_execs exec ON exec.comp_id = conn.provider_comp_id
	WHERE conn.conn_id = $1`

	selectSnapByChannel = `WITH acts AS (
		SELECT act.conn_id, act.act_nr, act.comp_id, act.comp_rn, act.chnl_id, act.kind, act.exp
		FROM pool_conn_acts act WHERE act.conn_id = $1 AND act.chnl_id = $2
	)
	SELECT conn.conn_id, conn.provider_comp_id, exec.comp_rn AS provider_comp_rn,
		conn.offset_nr, conn.act_nr,
		(SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM acts t WHERE t.act_nr > conn.offset_nr) AS acts
	FROM pool_comp_conns conn
	JOIN pool_comp_execs exec ON exec.comp_id = conn.provider_comp_id
	WHERE conn.conn_id = $1`
)
