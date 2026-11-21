package adapter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/compvar"
	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	compexec "orglang/go-engine/internal/proc/compexec/core"
	"orglang/go-engine/internal/proc/ent"
)

// Adapter
type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

// for compilation purposes
func newRepo() compexec.Repo {
	return new(daoEnt)
}

func newDaoEnt(client *ent.Client, log *slog.Logger) *daoEnt {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec compexec.ExecRec) error {
	dto := DataFromExecRec(rec)
	refAttr := slog.Any("ref", rec.CompRef)
	_, err := dao.client.CompExec.Create().
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

func (dao *daoEnt) GetSnapByRef(ctx context.Context, conn tx.Conn, ref compsem.SemRef) (compexec.ExecSnap, error) {
	refAttr := slog.Any("ref", ref)
	chnlRows, err := conn.Pgx.Query(ctx, selectChnls, ref.CompID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return compexec.ExecSnap{}, err
	}
	defer chnlRows.Close()
	chnlDtos, err := pgx.CollectRows(chnlRows, pgx.RowToStructByName[compvaradapter.VarRec])
	if err != nil {
		dao.log.Error("row scanning failed", refAttr, slog.Any("t", reflect.TypeOf(chnlDtos)))
		return compexec.ExecSnap{}, err
	}
	chnls, err := compvaradapter.DataToLinearRecs(chnlDtos)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return compexec.ExecSnap{}, err
	}
	dao.log.Debug("snap selection succeed", refAttr)
	return compexec.ExecSnap{
		LinearVars: compvar.IndexBy(compexec.ChnlPH, chnls),
	}, nil
}

func (dao *daoEnt) ModifyRec(ctx context.Context, conn tx.Conn, mod compexec.ExecMod) (err error) {
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
		ct, err := execRes.Exec()
		if err != nil {
			dao.log.Error("execution failed", slog.Any("dto", dto))
		}
		if ct.RowsAffected() == 0 {
			dao.log.Error("update failed")
			return errOptimisticUpdate(seqnum.ADT(dto.CompRN))
		}
	}
	if err != nil {
		return err
	}
	dao.log.Debug("update succeed")
	return nil
}

const (
	insertBind = `
		insert into proc_binds (
			impl_id, chnl_ph, chnl_id, state_id, impl_rn
		) values (
			@impl_id, @chnl_ph, @chnl_id, @state_id, @impl_rn
		)`

	updateExec = `
		update proc_comp_execs
		set comp_rn = @comp_rn + 1
		where comp_id = @comp_id
			and comp_rn = @comp_rn`

	selectChnls = `
		with bnds as not materialized (
			select distinct on (chnl_ph)
				*
			from proc_bnds
			where proc_id = 'proc1'
			order by chnl_ph, abs(rev) desc
		), liabs as not materialized (
			select distinct on (proc_id)
				*
			from pool_liabs
			where proc_id = 'proc1'
			order by proc_id, abs(rev) desc
		)
		select
			bnd.*,
			prvd.desc_id
		from bnds bnd
		left join liabs liab
			on liab.proc_id = bnd.proc_id
		left join pool_comp_execs prvd
			on prvd.desc_id = liab.desc_id`
)

func errOptimisticUpdate(got seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
