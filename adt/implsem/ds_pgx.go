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
	dto, err := DataFromRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}
	args := pgx.NamedArgs{
		"impl_id": dto.ImplID,
		"impl_rn": dto.ImplRN,
		"kind":    dto.Kind,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRec, args)
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}
	return nil
}

func (dao *pgxDAO) SelectRefsByQNs(source db.Source, implQNs []uniqsym.ADT) (_ map[uniqsym.ADT]SemRef, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection started", slog.Any("xactQNs", implQNs))
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
		qnAttr := slog.Any("xactQN", implQN)
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
	insertRec = `
		insert into impl_sems (
			impl_id, impl_rn, kind
		) values (
			@impl_id, @impl_rn, @kind
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
