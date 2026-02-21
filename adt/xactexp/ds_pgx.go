package xactexp

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/valkey"
)

// Adapter
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

func (dao *pgxDAO) InsertRec(source db.Source, rec ExpRec, ref descsem.SemRef) (err error) {
	ds := db.MustConform[db.SourcePgx](source)
	idAttr := slog.Any("expVK", rec.Key())
	dto := dataFromExpRec(rec)
	batch := pgx.Batch{}
	for _, st := range dto.States {
		args := pgx.NamedArgs{
			"exp_vk":     st.ExpVK,
			"sup_exp_vk": st.SupExpVK,
			"desc_id":    ref.DescID,
			"desc_rn":    ref.DescRN,
			"kind":       st.K,
			"spec":       st.Spec,
		}
		batch.Queue(insertRec, args)
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for range dto.States {
		_, err = br.Exec()
		if err != nil {
			dao.log.Error("query execution failed", idAttr)
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (dao *pgxDAO) SelectRecByVK(source db.Source, expVK valkey.ADT) (ExpRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	idAttr := slog.Any("expVK", expVK)
	rows, err := ds.Conn.Query(ds.Ctx, selectByID, valkey.ConvertToInteger(expVK))
	if err != nil {
		dao.log.Error("query execution failed", idAttr)
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[stateDS])
	if err != nil {
		dao.log.Error("row collection failed", idAttr)
		return nil, err
	}
	if len(dtos) == 0 { // revive:disable-line
		dao.log.Error("entity selection failed", idAttr)
		return nil, errors.New("no rows selected")
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entity selection succeed", slog.Any("dtos", dtos))
	states := make(map[int64]stateDS, len(dtos))
	for _, dto := range dtos {
		states[dto.ExpVK] = dto
	}
	return statesToExpRec(states, states[valkey.ConvertToInteger(expVK)])
}

func (dao *pgxDAO) SelectEnv(source db.Source, expVKs []valkey.ADT) (map[valkey.ADT]ExpRec, error) {
	recs, err := dao.SelectRecsByVKs(source, expVKs)
	if err != nil {
		return nil, err
	}
	env := make(map[valkey.ADT]ExpRec, len(recs))
	for _, rec := range recs {
		env[rec.Key()] = rec
	}
	return env, nil
}

func (dao *pgxDAO) SelectRecsByVKs(source db.Source, expVKs []valkey.ADT) (_ []ExpRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	batch := pgx.Batch{}
	for _, expVK := range expVKs {
		batch.Queue(selectByID, valkey.ConvertToInteger(expVK))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var recs []ExpRec
	for _, expVK := range expVKs {
		idAttr := slog.Any("expVK", expVK)
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", idAttr)
		}
		dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[stateDS])
		if err != nil {
			dao.log.Error("rows collection failed", idAttr)
		}
		if len(dtos) == 0 {
			dao.log.Error("entity selection failed", idAttr)
			return nil, ErrDoesNotExist(expVK)
		}
		rec, err := dataToExpRec(expRecDS{valkey.ConvertToInteger(expVK), dtos})
		if err != nil {
			dao.log.Error("model conversion failed", idAttr)
			return nil, err
		}
		recs = append(recs, rec)
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("recs", recs))
	return recs, err
}

const (
	insertRec = `
		insert into xact_exps (
			exp_vk, sup_exp_vk, desc_id, desc_rn, kind, spec
		) values (
			@exp_vk, @sup_exp_vk, @desc_id, @desc_rn, @kind, @spec
		)
		on conflict (exp_vk) do nothing`

	selectByID = `
		with recursive exp_tree AS (
			select top.*
			from xact_exps top
			where exp_vk = $1
			union all
			select sub.*
			from xact_exps sub, exp_tree sup
			where sub.sup_exp_vk = sup.exp_vk
		)
		select * from exp_tree`
)
