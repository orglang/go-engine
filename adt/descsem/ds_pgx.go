package descsem

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/adt/seqnum"
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

func (dao *pgxDAO) AddRec(source db.Source, rec SemRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	qnAttr := slog.Any("qn", rec.DescQN)
	dto, convertErr := DataFromRec(rec)
	if convertErr != nil {
		dao.log.Error("model conversion failed", qnAttr)
		return convertErr
	}
	args := pgx.NamedArgs{
		"desc_id": dto.DescID,
		"desc_rn": dto.DescRN,
		"desc_qn": dto.DescQN,
		"kind":    dto.Kind,
	}
	_, execErr1 := ds.Conn.Exec(ds.Ctx, insertRef, args)
	if execErr1 != nil {
		dao.log.Error("query execution failed", qnAttr)
		return execErr1
	}
	_, execErr2 := ds.Conn.Exec(ds.Ctx, insertBind, args)
	if execErr2 != nil {
		dao.log.Error("query execution failed", qnAttr)
		return execErr2
	}
	return nil
}

func (dao *pgxDAO) GetRefsByQNs(source db.Source, descQNs []uniqsym.ADT) (_ map[uniqsym.ADT]SemRef, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection started", slog.Any("xactQNs", descQNs))
	if len(descQNs) == 0 {
		return map[uniqsym.ADT]SemRef{}, nil
	}
	batch := pgx.Batch{}
	for _, descQN := range descQNs {
		batch.Queue(selectRefByQN, uniqsym.ConvertToString(descQN))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]SemRefDS, len(descQNs))
	for _, descQN := range descQNs {
		qnAttr := slog.Any("xactQN", descQN)
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
		dtos[descQN] = dto
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	return DataToRefMap(dtos)
}

const (
	insertRef = `
		insert into desc_sems (
			desc_id, desc_rn, kind
		) values (
			@desc_id, @desc_rn, @kind
		)`

	insertBind = `
		insert into desc_binds (
			desc_qn, desc_id
		) values (
			@desc_qn, @desc_id
		)`

	selectRefByQN = `
		select
			ds.desc_id,
			ds.desc_rn
		from desc_sems ds
		left join desc_binds db
			on db.desc_id = ds.desc_id
		where db.desc_qn = $1`
)

func errOptimisticUpdate(got seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
