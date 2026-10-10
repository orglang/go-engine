package adapter

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/typesem"
	"orglang/go-engine/adt/valkey"

	pooltypeexp "orglang/go-engine/pool/typeexp/core"
)

// Adapter
type daoPgx struct {
	qb  pooltypeexp.QueryBuilder
	log *slog.Logger
}

func NewDaoPgx(qb pooltypeexp.QueryBuilder, log *slog.Logger) pooltypeexp.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() pooltypeexp.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(uow db.UoW, rec pooltypeexp.ExpRec, ref typesem.SemRef) (err error) {
	vkAttr := slog.Any("vk", rec.Key())
	dto := dataFromExpRec(rec)
	batch := pgx.Batch{}
	for _, st := range dto.States {
		sql, args := dao.qb.InsertRec(st)
		batch.Queue(sql, args...)
	}
	br := uow.Pgx.SendBatch(uow.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for range dto.States {
		_, readErr := br.Exec()
		if readErr != nil {
			dao.log.Error("query execution failed", vkAttr)
			return readErr
		}
	}
	return nil
}

func (dao *daoPgx) GetRecByVK(uow db.UoW, expVK valkey.ADT) (pooltypeexp.ExpRec, error) {
	vkAttr := slog.Any("vk", expVK)
	sql := dao.qb.SelectRecByVK()
	rows, err := uow.Pgx.Query(uow.Ctx, sql, valkey.ConvertToInt(expVK))
	if err != nil {
		dao.log.Error("query execution failed", vkAttr, slog.String("sql", sql))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[pooltypeexp.StateDS])
	if err != nil {
		dao.log.Error("rows scanning failed", vkAttr)
		return nil, err
	}
	if len(dtos) == 0 { // revive:disable-line
		dao.log.Error("selection failed", vkAttr)
		return nil, errors.New("no rows selected")
	}
	dao.log.Log(uow.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	states := make(map[int64]pooltypeexp.StateDS, len(dtos))
	for _, dto := range dtos {
		states[dto.ExpVK] = dto
	}
	return statesToExpRec(states, states[valkey.ConvertToInt(expVK)])
}

func (dao *daoPgx) GetRecsByVKs(uow db.UoW, expVKs []valkey.ADT) (_ []pooltypeexp.ExpRec, err error) {
	batch := pgx.Batch{}
	sql := dao.qb.SelectRecByVK()
	for _, expVK := range expVKs {
		batch.Queue(sql, valkey.ConvertToInt(expVK))
	}
	br := uow.Pgx.SendBatch(uow.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	recs := make([]pooltypeexp.ExpRec, 0, len(expVKs))
	for _, expVK := range expVKs {
		vkAttr := slog.Any("vk", expVK)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", vkAttr, slog.String("sql", sql))
			return nil, readErr
		}
		dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[pooltypeexp.StateDS])
		if scanErr != nil {
			dao.log.Error("rows scanning failed", vkAttr)
			return nil, scanErr
		}
		if len(dtos) == 0 {
			dao.log.Error("selection failed", vkAttr)
			return nil, pooltypeexp.ErrDoesNotExist(expVK)
		}
		rec, convErr := dataToExpRec(pooltypeexp.ExpRecDS{ExpVK: valkey.ConvertToInt(expVK), States: dtos})
		if convErr != nil {
			dao.log.Error("model conversion failed", vkAttr)
			return nil, convErr
		}
		recs = append(recs, rec)
	}
	dao.log.Log(uow.Ctx, lf.LevelTrace, "selection succeed", slog.Any("recs", recs))
	return recs, err
}

func (dao *daoPgx) GetRecMap(uow db.UoW, expVKs map[symbol.ADT]valkey.ADT) (_ map[symbol.ADT]pooltypeexp.ExpRec, err error) {
	batch := pgx.Batch{}
	sql := dao.qb.SelectRecByVK()
	for _, expVK := range expVKs {
		batch.Queue(sql, valkey.ConvertToInt(expVK))
	}
	br := uow.Pgx.SendBatch(uow.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	recs := make(map[symbol.ADT]pooltypeexp.ExpRec, len(expVKs))
	for expPH, expVK := range expVKs {
		vkAttr := slog.Any("vk", expVK)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", vkAttr, slog.String("sql", sql))
			return nil, readErr
		}
		dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[pooltypeexp.StateDS])
		if scanErr != nil {
			dao.log.Error("rows scanning failed", vkAttr)
			return nil, scanErr
		}
		dao.log.Log(uow.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
		if len(dtos) == 0 {
			dao.log.Error("selection failed", vkAttr)
			return nil, pooltypeexp.ErrDoesNotExist(expVK)
		}
		rec, convErr := dataToExpRec(pooltypeexp.ExpRecDS{ExpVK: valkey.ConvertToInt(expVK), States: dtos})
		if convErr != nil {
			dao.log.Error("model conversion failed", vkAttr)
			return nil, convErr
		}
		recs[expPH] = rec
	}
	dao.log.Log(uow.Ctx, lf.LevelTrace, "getting succeed", slog.Any("recs", recs))
	return recs, err
}
