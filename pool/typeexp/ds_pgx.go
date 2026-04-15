package typeexp

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/semtype"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/valkey"
)

// Adapter
type pgxDAO struct {
	qb  queryBuilder
	log *slog.Logger
}

func newPgxDAO(qb queryBuilder, log *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{qb, log.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) AddRec(source db.Source, rec ExpRec, ref semtype.TypeRef) (err error) {
	ds := db.MustConform[db.SourcePgx](source)
	vkAttr := slog.Any("vk", rec.Key())
	dto := dataFromExpRec(rec)
	batch := pgx.Batch{}
	for _, st := range dto.States {
		sql, args := dao.qb.insertRec(st)
		batch.Queue(sql, args...)
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
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

func (dao *pgxDAO) GetRecByVK(source db.Source, expVK valkey.ADT) (ExpRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	vkAttr := slog.Any("vk", expVK)
	sql := dao.qb.selectRecByVK()
	rows, err := ds.Conn.Query(ds.Ctx, sql, valkey.ConvertToInt(expVK))
	if err != nil {
		dao.log.Error("query execution failed", vkAttr, slog.String("sql", sql))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[stateDS])
	if err != nil {
		dao.log.Error("rows scanning failed", vkAttr)
		return nil, err
	}
	if len(dtos) == 0 { // revive:disable-line
		dao.log.Error("selection failed", vkAttr)
		return nil, errors.New("no rows selected")
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	states := make(map[int64]stateDS, len(dtos))
	for _, dto := range dtos {
		states[dto.ExpVK] = dto
	}
	return statesToExpRec(states, states[valkey.ConvertToInt(expVK)])
}

func (dao *pgxDAO) GetRecsByVKs(source db.Source, expVKs []valkey.ADT) (_ []ExpRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	batch := pgx.Batch{}
	sql := dao.qb.selectRecByVK()
	for _, expVK := range expVKs {
		batch.Queue(sql, valkey.ConvertToInt(expVK))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	recs := make([]ExpRec, 0, len(expVKs))
	for _, expVK := range expVKs {
		vkAttr := slog.Any("vk", expVK)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", vkAttr, slog.String("sql", sql))
			return nil, readErr
		}
		dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[stateDS])
		if scanErr != nil {
			dao.log.Error("rows scanning failed", vkAttr)
			return nil, scanErr
		}
		if len(dtos) == 0 {
			dao.log.Error("selection failed", vkAttr)
			return nil, ErrDoesNotExist(expVK)
		}
		rec, convErr := dataToExpRec(expRecDS{valkey.ConvertToInt(expVK), dtos})
		if convErr != nil {
			dao.log.Error("model conversion failed", vkAttr)
			return nil, convErr
		}
		recs = append(recs, rec)
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("recs", recs))
	return recs, err
}

func (dao *pgxDAO) GetRecMap(source db.Source, expVKs map[symbol.ADT]valkey.ADT) (_ map[symbol.ADT]ExpRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	batch := pgx.Batch{}
	sql := dao.qb.selectRecByVK()
	for _, expVK := range expVKs {
		batch.Queue(sql, valkey.ConvertToInt(expVK))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	recs := make(map[symbol.ADT]ExpRec, len(expVKs))
	for expPH, expVK := range expVKs {
		vkAttr := slog.Any("vk", expVK)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", vkAttr, slog.String("sql", sql))
			return nil, readErr
		}
		dtos, scanErr := pgx.CollectRows(rows, pgx.RowToStructByName[stateDS])
		if scanErr != nil {
			dao.log.Error("rows scanning failed", vkAttr)
			return nil, scanErr
		}
		dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
		if len(dtos) == 0 {
			dao.log.Error("selection failed", vkAttr)
			return nil, ErrDoesNotExist(expVK)
		}
		rec, convErr := dataToExpRec(expRecDS{valkey.ConvertToInt(expVK), dtos})
		if convErr != nil {
			dao.log.Error("model conversion failed", vkAttr)
			return nil, convErr
		}
		recs[expPH] = rec
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting succeed", slog.Any("recs", recs))
	return recs, err
}
