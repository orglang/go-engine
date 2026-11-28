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
	commactadapter "orglang/go-engine/internal/proc/commact/adapter"
	commact "orglang/go-engine/internal/proc/commact/core"
	"orglang/go-engine/internal/proc/ent"
	entexeccomm "orglang/go-engine/internal/proc/ent/execcomm"
	"orglang/go-engine/internal/proc/execcomm/core"
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

func (dao *daoEnt) ModifyRec(ctx context.Context, conn tx.Conn, mod core.CommMod) error {
	refAttr := slog.Any("comm_id", mod.CommID)
	var offsetNr *int64
	if mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty() {
		value := int64(mod.OffsetNr.Get())
		offsetNr = &value
	}
	var actNr *int64
	if len(mod.Acts) > 0 {
		value := int64(mod.ProviderExecRN) + 1
		actNr = &value
	}

	if offsetNr != nil || actNr != nil {
		update := dao.client.ExecComm.Update().
			Where(entexeccomm.IDEQ(identity.ConvertToString(mod.CommID)))
		if offsetNr != nil {
			update.SetOffsetNr(*offsetNr)
		}
		if actNr != nil {
			update.SetActNr(*actNr)
		}
		affected, err := update.Save(entpgx.WithConn(ctx, conn.Pgx))
		if err != nil {
			dao.log.Error("Communication update failed", refAttr)
			return err
		}
		if affected == 0 {
			return fmt.Errorf("process communication missing: %v", mod.CommID)
		}
	}

	for _, act := range mod.Acts {
		dto, err := commactadapter.DataFromStepRec(act)
		if err != nil {
			return err
		}
		dto.ActNr = int64(mod.ProviderExecRN) + 1
		const insertAct = `
INSERT INTO proc_comm_acts (comm_id, act_nr, exec_id, exec_rn, chnl_id, kind, expr)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`
		_, err = conn.Pgx.Exec(ctx, insertAct, dto.CommID, dto.ActNr, dto.ExecID, dto.ExecRN, dto.ChnlID, dto.K, dto.Expr)
		if err != nil {
			dao.log.Error("Act insertion failed", slog.Any("act", act))
			return err
		}
	}
	return nil
}

func (dao *daoEnt) GetSnapByQry(ctx context.Context, conn tx.Conn, qry core.CommQry) (core.CommSnap, error) {
	commID := identity.ConvertToString(qry.CommID)
	query := `
WITH acts AS (
    SELECT act.*
    FROM proc_comm_acts act
    WHERE act.comm_id = $1
`
	args := []any{commID}
	if qry.ChnlID != nil && !qry.ChnlID.IsEmpty() {
		query += " AND act.chnl_id = $2"
		args = append(args, identity.ConvertToString(qry.ChnlID.Get()))
	}
	query += `
)
SELECT
    conn.comm_id,
    exec.exec_id AS provider_exec_id,
    exec.exec_rn AS provider_exec_rn,
    conn.offset_nr,
    conn.act_nr,
    (SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM acts t WHERE t.act_nr > conn.offset_nr) AS acts
FROM proc_exec_comms conn
JOIN proc_exec_evals exec ON exec.exec_id = conn.provider_exec_id
WHERE conn.comm_id = $1
`
	rows, err := conn.Pgx.Query(ctx, query, args...)
	if err != nil {
		dao.log.Error("Communication snapshot query failed", slog.Any("comm_id", qry.CommID))
		return core.CommSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[CommSnap])
	if errors.Is(err, pgx.ErrNoRows) {
		return core.CommSnap{CommID: qry.CommID}, nil
	}
	if err != nil {
		dao.log.Error("Communication snapshot scan failed", slog.Any("comm_id", qry.CommID))
		return core.CommSnap{}, err
	}
	parsedCommID, err := identity.ConvertFromString(dto.CommID)
	if err != nil {
		return core.CommSnap{}, err
	}
	providerCompID, err := identity.ConvertFromString(dto.ProviderExecID)
	if err != nil {
		return core.CommSnap{}, err
	}
	snap := core.CommSnap{
		CommID: parsedCommID, ProviderExecID: providerCompID,
		ProviderExecRN: seqnum.ADT(dto.ProviderExecRN), OffsetNr: seqnum.ADT(dto.OffsetNr),
		Acts: make([]commact.ActRec, 0, len(dto.Acts)),
	}
	for _, actDTO := range dto.Acts {
		act, err := commactadapter.DataToStepRec(actDTO)
		if err != nil {
			return core.CommSnap{}, err
		}
		snap.Acts = append(snap.Acts, act)
	}
	dao.log.Log(ctx, lf.LevelTrace, "Communication snapshot loaded", slog.Any("comm_id", qry.CommID))
	return snap, nil
}
