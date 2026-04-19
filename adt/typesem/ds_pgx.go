package typesem

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/uniqsym"
)

type pgxDAO struct {
	qb  queryBuilder
	log *slog.Logger
}

func NewPgxDAO(typeTable, descTable string) func(log *slog.Logger) *pgxDAO {
	return func(log *slog.Logger) *pgxDAO {
		name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
		return &pgxDAO{newSQLBuilder(typeTable, descTable), log.With(name)}
	}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) TouchRef(source db.Source, ref SemRef) error {
	ds := db.MustConform[db.SourcePgx](source)
	dto := DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	sql, args := dao.qb.updateRef(dto)
	ct, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	if ct.RowsAffected() == 0 {
		dao.log.Error("touching failed", refAttr)
		return errConcurrentModification(ref)
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "touching succeed", refAttr)
	return nil
}

func (dao *pgxDAO) GetRefsByQNs(source db.Source, typeQNs []uniqsym.ADT) (_ map[uniqsym.ADT]SemRef, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting started", slog.Any("qns", typeQNs))
	if len(typeQNs) == 0 {
		return map[uniqsym.ADT]SemRef{}, nil
	}
	batch := pgx.Batch{}
	for _, typeQN := range typeQNs {
		sql := dao.qb.selectRefByQN()
		batch.Queue(sql, uniqsym.ConvertToString(typeQN))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]SemRefDS, len(typeQNs))
	for _, typeQN := range typeQNs {
		qnAttr := slog.Any("qn", typeQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr)
			return nil, readErr
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[SemRefDS])
		if scanErr != nil {
			dao.log.Error("row scanning failed", qnAttr)
			return nil, scanErr
		}
		dtos[typeQN] = dto
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "getting succeed", slog.Any("dtos", dtos))
	return DataToRefMap(dtos)
}

func errConcurrentModification(got SemRef) error {
	return fmt.Errorf("concurrent modification: %v", got)
}
