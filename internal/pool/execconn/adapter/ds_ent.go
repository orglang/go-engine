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
	entexecconn "orglang/go-engine/internal/pool/ent/execconn"

	connactadapter "orglang/go-engine/internal/pool/connact/adapter"
	execconn "orglang/go-engine/internal/pool/execconn/core"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) execconn.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() execconn.Repo { return new(daoEnt) }

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec execconn.ConnRec) error {
	dto := DataFromRec(rec)
	refAttr := slog.Any("conn_id", rec.ConnID)
	txCtx := entpgx.WithConn(ctx, conn.Pgx)
	_, err := dao.client.ExecConn.Create().
		SetID(dto.ConnID).
		SetProviderExecID(dto.ProviderExecID).
		SetProviderExecRn(dto.ProviderExecRN).
		SetClientExecID(dto.ClientExecID).
		SetClientExecRn(dto.ClientExecRN).
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

func (dao *daoEnt) ModifyRec(ctx context.Context, conn tx.Conn, mod execconn.ConnMod) error {
	refAttr := slog.Any("conn_id", mod.ConnID)
	if (mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty()) || len(mod.Acts) > 0 {
		dto := DataFromMod(mod)
		if len(mod.Acts) > 0 {
			// The provider CAS advances exec_rn exactly once for this operation.
			// Mirror that revision as the Act number; Take appends at most one Act per operation.
			dto.ActNr = sql.Null[int64]{V: int64(mod.ProviderExecRN) + 1, Valid: true}
		}

		txCtx := entpgx.WithConn(ctx, conn.Pgx)
		update := dao.client.ExecConn.Update().
			Where(entexecconn.IDEQ(dto.ConnID))
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
			return execconn.ErrMissingConn(mod.ConnID)
		}
	}

	for _, act := range mod.Acts {
		dto := connactadapter.DataFromStepRec(act)
		dto.ActNr = int64(mod.ProviderExecRN) + 1
		_, err := conn.Pgx.Exec(ctx, insertActRec,
			dto.ConnID, dto.ActNr, dto.ExecID, dto.ExecRN, dto.ChnlID, dto.K, dto.Exp,
		)
		if err != nil {
			dao.log.Error("Act insertion failed", refAttr)
			return err
		}
	}

	dao.log.Log(ctx, lf.LevelTrace, "Connection modification succeeded", refAttr)
	return nil
}

func (dao *daoEnt) GetSnapByQry(ctx context.Context, conn tx.Conn, qry execconn.ConnQry) (execconn.ConnSnap, error) {
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
		return execconn.ConnSnap{}, err
	}
	defer rows.Close()

	snapDTO, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ConnSnap])
	if err != nil {
		dao.log.Error("Connection snapshot scan failed", refAttr)
		return execconn.ConnSnap{}, err
	}
	snap, err := DataToSnap(snapDTO)
	if err != nil {
		dao.log.Error("Connection snapshot conversion failed", refAttr)
		return execconn.ConnSnap{}, err
	}
	providerCompID, err := identity.ConvertFromString(snapDTO.ProviderExecID)
	if err != nil {
		return execconn.ConnSnap{}, err
	}
	snap.ProviderExecID = providerCompID
	snap.ProviderExecRN = seqnum.ADT(snapDTO.ProviderExecRN)
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot loaded", slog.Any("dto", snapDTO))
	return snap, nil
}

const (
	insertActRec = `
INSERT INTO pool_conn_acts
(conn_id, act_nr, exec_id, exec_rn, chnl_id, kind, exp)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

	selectSnap = `
WITH acts AS (
    SELECT act.conn_id, act.act_nr, act.exec_id, act.exec_rn, act.chnl_id, act.kind, act.exp
    FROM pool_conn_acts act WHERE act.conn_id = $1
)
SELECT conn.conn_id, conn.provider_exec_id, exec.exec_rn AS provider_exec_rn,
    conn.offset_nr, conn.act_nr,
    (SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM acts t WHERE t.act_nr > conn.offset_nr) AS acts
FROM pool_exec_conns conn
JOIN pool_exec_evals exec ON exec.exec_id = conn.provider_exec_id
WHERE conn.conn_id = $1
`

	selectSnapByChannel = `
WITH acts AS (
    SELECT act.conn_id, act.act_nr, act.exec_id, act.exec_rn, act.chnl_id, act.kind, act.exp
    FROM pool_conn_acts act WHERE act.conn_id = $1 AND act.chnl_id = $2
)
SELECT conn.conn_id, conn.provider_exec_id, exec.exec_rn AS provider_exec_rn,
    conn.offset_nr, conn.act_nr,
    (SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM acts t WHERE t.act_nr > conn.offset_nr) AS acts
FROM pool_exec_conns conn
JOIN pool_exec_evals exec ON exec.exec_id = conn.provider_exec_id
WHERE conn.conn_id = $1
`
)
