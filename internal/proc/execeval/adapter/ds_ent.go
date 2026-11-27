package adapter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/compvar"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	execeval "orglang/go-engine/internal/proc/execeval/core"
	"orglang/go-engine/internal/proc/ent"
)

// Adapter
type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

// for compilation purposes
func newRepo() execeval.Repo {
	return new(daoEnt)
}

func newDaoEnt(client *ent.Client, log *slog.Logger) *daoEnt {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec execeval.EvalRec) error {
	dto := DataFromEvalRec(rec)
	refAttr := slog.Any("ref", rec.CompRef)
	_, err := dao.client.ExecEval.Create().
		SetID(dto.CompID).
		SetCompRn(dto.CompRN).
		SetLiabMode(dto.LiabMode).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "insertion succeed", slog.Any("dto", dto))
	return nil
}

func (dao *daoEnt) GetSnapByRef(ctx context.Context, conn tx.Conn, ref comproot.SemRef) (execeval.EvalSnap, error) {
	refAttr := slog.Any("ref", ref)
	chnlRows, err := conn.Pgx.Query(ctx, selectChnls, ref.CompID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return execeval.EvalSnap{}, err
	}
	defer chnlRows.Close()
	chnlDtos, err := pgx.CollectRows(chnlRows, pgx.RowToStructByName[compvaradapter.VarRec])
	if err != nil {
		dao.log.Error("row scanning failed", refAttr, slog.Any("t", reflect.TypeOf(chnlDtos)))
		return execeval.EvalSnap{}, err
	}
	chnls, err := compvaradapter.DataToLinearRecs(chnlDtos)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return execeval.EvalSnap{}, err
	}
	dao.log.Debug("snap selection succeed", refAttr)
	return execeval.EvalSnap{
		LinearVars: compvar.IndexBy(execeval.ChnlPH, chnls),
	}, nil
}

func (dao *daoEnt) ModifyRec(ctx context.Context, conn tx.Conn, mod execeval.EvalMod) (err error) {
	if len(mod.CompRefs) == 0 {
		panic("empty locks")
	}
	dto, err := DataFromMod(mod)
	if err != nil {
		dao.log.Error("conversion failed")
		return err
	}
	// binds
	bindReq := pgx.Batch{}
	for _, dto := range dto.LinearVars {
		args := pgx.NamedArgs{
			"impl_id":  dto.CompID,
			"chnl_ph":  dto.ChnlPH,
			"chnl_id":  dto.ChnlID,
			"state_id": dto.ExpCK,
		}
		bindReq.Queue(insertBind, args)
	}
	if bindReq.Len() > 0 {
		bindRes := conn.Pgx.SendBatch(ctx, &bindReq)
		defer func() {
			err = errors.Join(err, bindRes.Close())
		}()
		for _, dto := range dto.LinearVars {
			_, err = bindRes.Exec()
			if err != nil {
				dao.log.Error("execution failed", slog.Any("dto", dto))
			}
		}
		if err != nil {
			return err
		}
	}
	// execs
	execReq := pgx.Batch{}
	for _, dto := range dto.CompRefs {
		args := pgx.NamedArgs{
			"comp_id": dto.CompID,
			"comp_rn": dto.CompRN,
		}
		execReq.Queue(updateExec, args)
	}
	execRes := conn.Pgx.SendBatch(ctx, &execReq)
	defer func() {
		err = errors.Join(err, execRes.Close())
	}()
	for _, dto := range dto.CompRefs {
		ct, execErr := execRes.Exec()
		if execErr != nil {
			dao.log.Error("execution failed", slog.Any("dto", dto))
			return execErr
		}
		if ct.RowsAffected() == 0 {
			dao.log.Error("update failed", slog.Any("dto", dto))
			return errOptimisticUpdate(seqnum.ADT(dto.CompRN))
		}
	}
	dao.log.Debug("update succeed")
	return nil
}

const (
	insertBind = `
INSERT INTO proc_binds (
    impl_id, chnl_ph, chnl_id, state_id, impl_rn
) VALUES (
    @impl_id, @chnl_ph, @chnl_id, @state_id, @impl_rn
)
`

	updateExec = `
UPDATE proc_exec_evals
SET comp_rn = @comp_rn + 1
WHERE comp_id = @comp_id
    AND comp_rn = @comp_rn
`

	selectChnls = `
WITH bnds AS NOT MATERIALIZED (
    SELECT DISTINCT ON (chnl_ph)
        *
    FROM proc_bnds
    WHERE proc_id = 'proc1'
    ORDER BY chnl_ph, abs(rev) DESC
), liabs AS NOT MATERIALIZED (
    SELECT DISTINCT ON (proc_id)
        *
    FROM pool_liabs
    WHERE proc_id = 'proc1'
    ORDER BY proc_id, abs(rev) DESC
)
SELECT
    bnd.*,
    prvd.desc_id
FROM bnds bnd
LEFT JOIN liabs liab
    ON liab.proc_id = bnd.proc_id
LEFT JOIN pool_exec_evals prvd
    ON prvd.desc_id = liab.desc_id
`
)

func errOptimisticUpdate(got seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
