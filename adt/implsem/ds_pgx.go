package implsem

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/uniqsym"
)

type pgxDAO struct {
	log *slog.Logger
}

func newPgxDAO(l *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{l.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) InsertRec(source db.Source, rec SemRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", rec.Ref)
	dto, convertErr := DataFromRec(rec)
	if convertErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return convertErr
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

func (dao *pgxDAO) SelectRefsByQNs(source db.Source, implQNs []uniqsym.ADT) (_ map[uniqsym.ADT]SemRef, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "starting selection...", slog.Any("qns", implQNs))
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
		dto, collectErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[SemRefDS])
		if collectErr != nil {
			dao.log.Error("row collection failed", qnAttr)
			return nil, collectErr
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

func errOptimisticUpdate(got revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
