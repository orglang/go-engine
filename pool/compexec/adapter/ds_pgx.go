package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/lf"
	"orglang/go-engine/lib/tx"

	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/uniqsym"

	compexec "orglang/go-engine/pool/compexec/core"
)

type daoPgx struct {
	qb  compexec.QueryBuilder
	log *slog.Logger
}

func NewDaoPgx(qb compexec.QueryBuilder, log *slog.Logger) compexec.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() compexec.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec compexec.ExecRec) error {
	dto := DataFromExecRec(rec)
	refAttr := slog.Any("ref", rec.CompRef)
	sql, args := dao.qb.InsertRec(dto)
	_, execErr := conn.Pgx.Exec(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", slog.Any("dto", dto))
	return nil
}

func (dao *daoPgx) ModifyRec(context.Context, tx.Conn, compexec.ExecMod) error {
	panic("unimplemented")
}

func (dao *daoPgx) GetSnapByRef(ctx context.Context, conn tx.Conn, ref compsem.SemRef) (compexec.ExecSnap2, error) {
	refAttr := slog.Any("ref", ref)
	sql, args := dao.qb.SelectRecByRef(compsem.DataFromRef(ref))
	rows, execErr := conn.Pgx.Query(ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return compexec.ExecSnap2{}, execErr
	}
	defer rows.Close()
	dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[compexec.ExecSnap2DS])
	if scanErr != nil {
		dao.log.Error("rows scanning failed", refAttr, slog.Any("dto", dto))
		return compexec.ExecSnap2{}, scanErr
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dto", dto))
	rec, convErr := DataToExecSnap2(dto)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return compexec.ExecSnap2{}, convErr
	}
	return rec, nil
}

func (dao *daoPgx) GetSnapMapByQNs(ctx context.Context, conn tx.Conn, termQNs []uniqsym.ADT) (_ map[uniqsym.ADT]compexec.ExecSnap1, err error) {
	dao.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("qns", termQNs))
	if len(termQNs) == 0 {
		return map[uniqsym.ADT]compexec.ExecSnap1{}, nil
	}
	batch := pgx.Batch{}
	for _, termQN := range termQNs {
		sql, args := dao.qb.SelectSnapByQN(uniqsym.ConvertToString(termQN))
		batch.Queue(sql, args...)
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]compexec.ExecSnap1DS, len(termQNs))
	for _, termQN := range termQNs {
		qnAttr := slog.Any("qn", termQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr)
			return nil, readErr
		}
		var dto compexec.ExecSnap1DS
		scanErr := pgxscan.ScanOne(dto, rows)
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return nil, scanErr
		}
		dtos[termQN] = dto
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return DataToSnapMap(dtos)
}
