package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/comproot"
	compsemadapter "orglang/go-engine/internal/adt/comproot/adapter"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"

	compexec "orglang/go-engine/internal/pool/compexec/core"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) compexec.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() compexec.Repo { return new(daoEnt) }

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
	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", slog.Any("dto", dto))
	return nil
}

func (dao *daoEnt) ModifyRec(context.Context, tx.Conn, compexec.ExecMod) error {
	panic("unimplemented")
}

func (dao *daoEnt) GetSnapByRef(ctx context.Context, conn tx.Conn, ref comproot.SemRef) (compexec.ExecSnap2, error) {
	refAttr := slog.Any("ref", ref)
	rows, err := conn.Pgx.Query(ctx, selectSnapByRef,
		compsemadapter.DataFromRef(ref).CompID,
		compsemadapter.DataFromRef(ref).CompID,
		compsemadapter.DataFromRef(ref).CompID,
	)
	if err != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", selectSnapByRef))
		return compexec.ExecSnap2{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ExecSnap2])
	if err != nil {
		dao.log.Error("rows scanning failed", refAttr, slog.Any("dto", dto))
		return compexec.ExecSnap2{}, err
	}
	rec, err := DataToExecSnap2(dto)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return compexec.ExecSnap2{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dto", dto))
	return rec, nil
}

func (dao *daoEnt) GetSnapMapByQNs(ctx context.Context, conn tx.Conn, termQNs []uniqsym.ADT) (_ map[uniqsym.ADT]compexec.ExecSnap1, err error) {
	dao.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("qns", termQNs))
	if len(termQNs) == 0 {
		return map[uniqsym.ADT]compexec.ExecSnap1{}, nil
	}
	batch := pgx.Batch{}
	for _, termQN := range termQNs {
		batch.Queue(selectSnapByQN, uniqsym.ConvertToString(termQN))
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]ExecSnap1, len(termQNs))
	for _, termQN := range termQNs {
		qnAttr := slog.Any("qn", termQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr)
			return nil, readErr
		}
		var dto ExecSnap1
		scanErr := pgxscan.ScanOne(&dto, rows)
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return nil, scanErr
		}
		dtos[termQN] = dto
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return DataToSnapMap(dtos)
}

const (
	selectSnapByRef = `
WITH struct_vars AS (
    SELECT var.comp_id, var.comp_rn, var.conn_id, var.chnl_id, var.chnl_ph, var.exp_ck, var.side
    FROM pool_struct_vars var WHERE comp_id = $1
), linear_vars AS (
    SELECT var.comp_id, var.comp_rn, var.conn_id, var.chnl_id, var.chnl_ph, var.exp_ck, var.side
    FROM pool_linear_vars var WHERE comp_id = $2
)
SELECT exec.comp_id, exec.comp_rn, exec.liab_mode,
    (SELECT array_agg(row(r.*)) FROM struct_vars r) AS struct_vars,
    (SELECT array_agg(row(r.*)) FROM linear_vars r) AS linear_vars
FROM pool_comp_execs exec WHERE exec.comp_id = $3
`

	selectSnapByQN = `
SELECT sem.comp_id, sem.comp_rn, exec.liab_mode,
    struct_var.comp_id, struct_var.comp_rn, struct_var.conn_id, struct_var.chnl_id,
    struct_var.chnl_ph, struct_var.exp_ck, struct_var.side,
    linear_var.comp_id, linear_var.comp_rn, linear_var.conn_id, linear_var.chnl_id,
    linear_var.chnl_ph, linear_var.exp_ck, linear_var.side
FROM pool_comp_execs sem
JOIN pool_comp_execs exec ON exec.comp_id = sem.comp_id
JOIN pool_struct_vars struct_var ON struct_var.comp_id = sem.comp_id AND struct_var.side = 1 AND exec.liab_mode = 1
JOIN pool_linear_vars linear_var ON linear_var.comp_id = sem.comp_id AND linear_var.side = 1 AND exec.liab_mode = 2
JOIN pool_impl_binds bind ON bind.impl_id = sem.comp_id
WHERE bind.impl_qn = $1
ORDER BY struct_var.comp_rn DESC, linear_var.comp_rn DESC
LIMIT 1
`
)
