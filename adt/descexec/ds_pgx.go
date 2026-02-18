package descexec

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"
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

func (dao *pgxDAO) InsertRec(source db.Source, rec ExecRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", rec.Ref)
	dto, err := DataFromRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}
	args := pgx.NamedArgs{
		"desc_id": dto.DescID,
		"desc_rn": dto.DescRN,
		"kind":    dto.Kind,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRec, args)
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}
	return nil
}

func (dao *pgxDAO) SelectRefsByQNs(source db.Source, descQNs []uniqsym.ADT) (_ map[uniqsym.ADT]ExecRef, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection started", slog.Any("xactQNs", descQNs))
	if len(descQNs) == 0 {
		return map[uniqsym.ADT]ExecRef{}, nil
	}
	batch := pgx.Batch{}
	for _, descQN := range descQNs {
		batch.Queue(selectRefByQN, uniqsym.ConvertToString(descQN))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	dtos := make(map[uniqsym.ADT]ExecRefDS, len(descQNs))
	for _, descQN := range descQNs {
		qnAttr := slog.Any("xactQN", descQN)
		rows, readErr := br.Query()
		if readErr != nil {
			dao.log.Error("query execution failed", qnAttr)
			return nil, readErr
		}
		dto, collectErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[ExecRefDS])
		if collectErr != nil {
			dao.log.Error("row collection failed", qnAttr)
			return nil, collectErr
		}
		dtos[descQN] = dto
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	return DataToDescRefs(dtos)
}

const (
	insertRec = `
		insert into desc_execs (
			desc_id, desc_rn, kind
		) values (
			@desc_id, @desc_rn, @kind
		)`

	touchRec = `
		update desc_binds
		set ref_rn = ref_rn + 1
		where desc_qn = @desc_qn
		and ref_rn = @ref_rn
	`

	selectRefByQN = `
		select
			de.desc_id,
			de.desc_rn
		from desc_execs de
		left join desc_binds db
			on db.desc_id = de.desc_id
		where db.desc_qn = $1`
)

func errOptimisticUpdate(got revnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}
