package implsem

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/seqnum"
	"orglang/go-engine/adt/uniqsym"
)

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

func (dao *pgxDAO) AddRec(source db.Source, rec SemRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", rec.ImplRef)
	dto, convErr := DataFromRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return convErr
	}
	args := pgx.NamedArgs{
		"impl_id": dto.ImplID,
		"impl_rn": dto.ImplRN,
		"impl_qn": dto.ImplQN,
		"kind":    dto.Kind,
	}
	_, execErr1 := ds.Conn.Exec(ds.Ctx, insertRef, args)
	if execErr1 != nil {
		dao.log.Error("query execution failed", refAttr)
		return execErr1
	}
	if !dto.ImplQN.Valid {
		return nil
	}
	_, execErr2 := ds.Conn.Exec(ds.Ctx, insertBind, args)
	if execErr2 != nil {
		dao.log.Error("query execution failed", refAttr)
		return execErr2
	}
	return nil
}

func (dao *pgxDAO) TouchRec(source db.Source, ref SemRef) error {
	ds := db.MustConform[db.SourcePgx](source)
	dto := DataFromRef(ref)
	refAttr := slog.Any("ref", ref)
	sql, args := dao.qb.updateRec(dto)
	ct, execErr := ds.Conn.Exec(ds.Ctx, sql, args...)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr, slog.String("sql", sql))
		return execErr
	}
	if ct.RowsAffected() == 0 {
		dao.log.Error("touch failed", refAttr)
		return ErrConcurrentModification(ref)
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "touch succeed", slog.Any("dto", dto))
	return nil
}

func (dao *pgxDAO) GetRefsByQNs(source db.Source, implQNs []uniqsym.ADT) (_ map[uniqsym.ADT]SemRef, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection started", slog.Any("qns", implQNs))
	if len(implQNs) == 0 {
		return map[uniqsym.ADT]SemRef{}, nil
	}
	batch := pgx.Batch{}
	for _, implQN := range implQNs {
		batch.Queue(selectRefByQN, uniqsym.ConvertToString(implQN))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]SemRefDS, len(implQNs))
	for _, implQN := range implQNs {
		qnAttr := slog.Any("qn", implQN)
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
		dtos[implQN] = dto
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	return DataToRefMap(dtos)
}

const (
	insertRef = `
		insert into impl_sems (
			impl_id, impl_rn, kind
		) values (
			@impl_id, @impl_rn, @kind
		)`

	insertBind = `
		insert into impl_binds (
			impl_qn, impl_id
		) values (
			@impl_qn, @impl_id
		)`

	selectRefByQN = `
		select
			is.impl_id,
			is.impl_rn
		from impl_sems is
		left join impl_binds ib
			on ib.impl_id = is.impl_id
		where ib.impl_qn = $1`
)

func errOptimisticUpdate(got seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
