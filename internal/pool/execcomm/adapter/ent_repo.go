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
	entexeccomm "orglang/go-engine/internal/pool/ent/execcomm"

	commactadapter "orglang/go-engine/internal/pool/commact/adapter"
	execcomm "orglang/go-engine/internal/pool/execcomm/core"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) execcomm.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() execcomm.Repo { return new(daoEnt) }

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec execcomm.CommRec) error {
	dto := DataFromRec(rec)
	refAttr := slog.Any("comm_id", rec.CommID)
	txCtx := entpgx.WithConn(ctx, conn.Pgx)
	_, err := dao.client.ExecComm.Create().
		SetID(dto.CommID).
		SetProviderExecID(dto.ProviderExecID).
		SetProviderExecRn(dto.ProviderExecRN).
		SetClientExecID(dto.ClientExecID).
		SetClientExecRn(dto.ClientExecRN).
		SetOffsetNr(dto.OffsetNr).
		SetActNr(dto.ActNr).
		Save(txCtx)
	if err != nil {
		dao.log.Error("Communication insertion failed", refAttr)
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "Communication insertion succeeded", refAttr)
	return nil
}

func (dao *daoEnt) ModifyRec(ctx context.Context, conn tx.Conn, mod execcomm.CommMod) error {
	refAttr := slog.Any("comm_id", mod.CommID)
	if (mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty()) || len(mod.Acts) > 0 {
		dto := DataFromMod(mod)
		if len(mod.Acts) > 0 {
			// The provider CAS advances exec_rn exactly once for this operation.
			// Mirror that revision as the Act number; Take appends at most one Act per operation.
			dto.ActNr = sql.Null[int64]{V: int64(mod.ProviderExecRN) + 1, Valid: true}
		}

		txCtx := entpgx.WithConn(ctx, conn.Pgx)
		update := dao.client.ExecComm.Update().
			Where(entexeccomm.IDEQ(dto.CommID))
		if dto.OffsetNr.Valid {
			update.SetOffsetNr(dto.OffsetNr.V)
		}
		if dto.ActNr.Valid {
			update.SetActNr(dto.ActNr.V)
		}
		affected, err := update.Save(txCtx)
		if err != nil {
			dao.log.Error("Communication offset update failed", refAttr)
			return err
		}
		if affected == 0 {
			return execcomm.ErrMissingComm(mod.CommID)
		}
	}

	for _, act := range mod.Acts {
		dto := commactadapter.DataFromStepRec(act)
		dto.ActNr = int64(mod.ProviderExecRN) + 1
		_, err := conn.Pgx.Exec(ctx, insertActRec,
			dto.CommID, dto.ActNr, dto.ExecID, dto.ExecRN, dto.ChnlID, dto.K, dto.Expr,
		)
		if err != nil {
			dao.log.Error("Act insertion failed", refAttr)
			return err
		}
	}

	dao.log.Log(ctx, lf.LevelTrace, "Communication modification succeeded", refAttr)
	return nil
}

func (dao *daoEnt) GetSnapByQry(ctx context.Context, conn tx.Conn, qry execcomm.CommQry) (execcomm.CommSnap, error) {
	refAttr := slog.Any("comm_id", qry.CommID)
	qryDTO := DataFromQry(qry)
	querySQL := selectSnap
	args := []any{qryDTO.CommID}
	if qryDTO.ChnlID.Valid {
		querySQL = selectSnapByChannel
		args = append(args, qryDTO.ChnlID.V)
	}
	dao.log.Log(ctx, lf.LevelTrace, "Communication snapshot selection started", slog.Any("qry", qryDTO))
	rows, err := conn.Pgx.Query(ctx, querySQL, args...)
	if err != nil {
		dao.log.Error("Communication snapshot query failed", refAttr)
		return execcomm.CommSnap{}, err
	}
	defer rows.Close()

	snapDTO, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[CommSnap])
	if err != nil {
		dao.log.Error("Communication snapshot scan failed", refAttr)
		return execcomm.CommSnap{}, err
	}
	snap, err := DataToSnap(snapDTO)
	if err != nil {
		dao.log.Error("Communication snapshot conversion failed", refAttr)
		return execcomm.CommSnap{}, err
	}
	providerCompID, err := identity.ConvertFromString(snapDTO.ProviderExecID)
	if err != nil {
		return execcomm.CommSnap{}, err
	}
	snap.ProviderExecID = providerCompID
	snap.ProviderExecRN = seqnum.ADT(snapDTO.ProviderExecRN)
	dao.log.Log(ctx, lf.LevelTrace, "Communication snapshot loaded", slog.Any("dto", snapDTO))
	return snap, nil
}

const (
	insertActRec = `
INSERT INTO pool_comm_acts
(comm_id, act_nr, exec_id, exec_rn, chnl_id, kind, expr)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`

	selectSnap = `
WITH acts AS (
    SELECT act.comm_id, act.act_nr, act.exec_id, act.exec_rn, act.chnl_id, act.kind, act.expr
    FROM pool_comm_acts act WHERE act.comm_id = $1
)
SELECT conn.comm_id, conn.provider_exec_id, exec.exec_rn AS provider_exec_rn,
    conn.offset_nr, conn.act_nr,
    (SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM acts t WHERE t.act_nr > conn.offset_nr) AS acts
FROM pool_exec_comms conn
JOIN pool_exec_evals exec ON exec.exec_id = conn.provider_exec_id
WHERE conn.comm_id = $1
`

	selectSnapByChannel = `
WITH acts AS (
    SELECT act.comm_id, act.act_nr, act.exec_id, act.exec_rn, act.chnl_id, act.kind, act.expr
    FROM pool_comm_acts act WHERE act.comm_id = $1 AND act.chnl_id = $2
)
SELECT conn.comm_id, conn.provider_exec_id, exec.exec_rn AS provider_exec_rn,
    conn.offset_nr, conn.act_nr,
    (SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM acts t WHERE t.act_nr > conn.offset_nr) AS acts
FROM pool_exec_comms conn
JOIN pool_exec_evals exec ON exec.exec_id = conn.provider_exec_id
WHERE conn.comm_id = $1
`
)
