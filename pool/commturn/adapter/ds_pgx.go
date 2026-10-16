package adapter

import (
	"errors"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"github.com/jackc/pgx/v5"

	commturn "orglang/go-engine/pool/commturn/core"
)

type daoPgx struct {
	qb  commturn.QueryBuilder
	log *slog.Logger
}

func NewDaoPgx(qb commturn.QueryBuilder, log *slog.Logger) commturn.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() commturn.Repo {
	return new(daoPgx)
}

func (dao *daoPgx) AddRec(db.UoW, commturn.TurnRec) error {
	panic("unimplemented")
}

func (dao *daoPgx) AddRecs(uow db.UoW, recs []commturn.TurnRec) (err error) {
	dao.log.Log(uow.Ctx, lf.LevelTrace, "insertion started", slog.Any("recs", recs))
	batch := pgx.Batch{}
	for _, rec := range recs {
		dto := DataFromStepRec(rec)
		sql, args := dao.qb.InsertRec(dto)
		batch.Queue(sql, args...)
	}
	br := uow.Pgx.SendBatch(uow.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for _, rec := range recs {
		_, readErr := br.Exec()
		if readErr != nil {
			dao.log.Error("query execution failed", slog.Any("rec", rec))
			return readErr
		}
	}
	dao.log.Log(uow.Ctx, lf.LevelTrace, "insertion succeed")
	return nil
}
