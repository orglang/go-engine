package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	"github.com/jackc/pgx/v5"

	commturn "orglang/go-engine/internal/pool/commturn/core"
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

func (dao *daoPgx) AddRec(context.Context, tx.Conn, commturn.TurnRec) error {
	panic("unimplemented")
}

func (dao *daoPgx) AddRecs(ctx context.Context, conn tx.Conn, recs []commturn.TurnRec) (err error) {
	dao.log.Log(ctx, lf.LevelTrace, "insertion started", slog.Any("recs", recs))
	batch := pgx.Batch{}
	for _, rec := range recs {
		dto := DataFromStepRec(rec)
		sql, args := dao.qb.InsertRec(dto)
		batch.Queue(sql, args...)
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
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
	dao.log.Log(ctx, lf.LevelTrace, "insertion succeed")
	return nil
}
