package typeexp

import (
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"

	"orglang/go-engine/adt/identity"
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

func (dao *pgxDAO) InsertRec(source db.Source, rec ExpRec) (err error) {
	ds := db.MustConform[db.SourcePgx](source)
	idAttr := slog.Any("expID", rec.Ident())
	dto := dataFromExpRec(rec)
	batch := pgx.Batch{}
	for _, st := range dto.States {
		sa := pgx.NamedArgs{
			"exp_id":     st.ExpID,
			"kind":       st.K,
			"sup_exp_id": st.SupExpID,
			"spec":       st.Spec,
		}
		batch.Queue(insertRec, sa)
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for range dto.States {
		_, err = br.Exec()
		if err != nil {
			dao.log.Error("query execution failed", idAttr, slog.String("q", insertRec))
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (dao *pgxDAO) SelectRecByID(source db.Source, expID identity.ADT) (ExpRec, error) {
	ds := db.MustConform[db.SourcePgx](source)
	idAttr := slog.Any("expID", expID)
	query := `
		WITH RECURSIVE top_states AS (
			select te.*
			from type_exps te
			WHERE exp_id = $1
			UNION ALL
			select be.*
			from type_exps be, top_states ts
			WHERE be.sup_exp_id = ts.exp_id
		)
		select * from top_states`
	rows, err := ds.Conn.Query(ds.Ctx, query, expID.String())
	if err != nil {
		dao.log.Error("query execution failed", idAttr, slog.String("q", query))
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
	states := make(map[string]stateDS, len(dtos))
	for _, dto := range dtos {
		states[dto.ExpID] = dto
	}
	return statesToExpRec(states, states[expID.String()])
}

func (dao *pgxDAO) SelectEnv(source db.Source, expIDs []identity.ADT) (map[identity.ADT]ExpRec, error) {
	recs, err := dao.SelectRecsByIDs(source, expIDs)
	if err != nil {
		return nil, err
	}
	env := make(map[identity.ADT]ExpRec, len(recs))
	for _, rec := range recs {
		env[rec.Ident()] = rec
	}
	return env, nil
}

func (dao *pgxDAO) SelectRecsByIDs(source db.Source, expIDs []identity.ADT) (_ []ExpRec, err error) {
	ds := db.MustConform[db.SourcePgx](source)
	batch := pgx.Batch{}
	for _, expID := range expIDs {
		batch.Queue(selectByID, expID.String())
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var recs []ExpRec
	for _, expID := range expIDs {
		idAttr := slog.Any("expID", expID)
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", idAttr, slog.String("q", selectByID))
		}
		dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[stateDS])
		if err != nil {
			dao.log.Error("rows collection failed", idAttr)
		}
		if len(dtos) == 0 {
			dao.log.Error("entity selection failed", idAttr)
			return nil, ErrDoesNotExist(expID)
		}
		rec, err := dataToExpRec(expRecDS{expID.String(), dtos})
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
		insert into type_exps (
			exp_id, kind, sup_exp_id, spec
		) values (
			@exp_id, @kind, @sup_exp_id, @spec
		)`

	selectByID = `
		WITH RECURSIVE state_tree AS (
			select root.*
			from type_exps root
			WHERE id = $1
			UNION ALL
			select child.*
			from type_exps child, state_tree parent
			WHERE child.sup_exp_id = parent.id
		)
		select * from state_tree`
)
