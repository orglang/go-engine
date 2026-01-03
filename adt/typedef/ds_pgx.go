package typedef

import (
	"errors"
	"fmt"
	"log/slog"
	"math"

	"github.com/jackc/pgx/v5"

	"orglang/orglang/lib/lf"
	"orglang/orglang/lib/sd"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"
)

// Adapter
type daoPgx struct {
	log *slog.Logger
}

func newDaoPgx(l *slog.Logger) *daoPgx {
	return &daoPgx{l}
}

// for compilation purposes
func newRepo() Repo {
	return &daoPgx{}
}

func (d *daoPgx) InsertType(source sd.Source, rec DefRec) error {
	ds := sd.MustConform[sd.SourcePgx](source)
	idAttr := slog.Any("defID", rec.DefID)
	d.log.Log(ds.Ctx, lf.LevelTrace, "entity insertion started", idAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil {
		d.log.Error("model mapping failed", idAttr)
		return err
	}
	insertRoot := `
		insert into type_def_roots (
			def_id, def_rn, title
		) values (
			@def_id, @def_rn, @title
		)`
	rootArgs := pgx.NamedArgs{
		"def_id": dto.DefID,
		"def_rn": dto.DefRN,
		"title":  dto.Title,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRoot, rootArgs)
	if err != nil {
		d.log.Error("query execution failed", idAttr, slog.String("q", insertRoot))
		return err
	}
	insertState := `
		insert into type_def_states (
			def_id, term_id, from_rn, to_rn
		) values (
			@def_id, @term_id, @from_rn, @to_rn
		)`
	stateArgs := pgx.NamedArgs{
		"def_id":  dto.DefID,
		"from_rn": dto.DefRN,
		"to_rn":   math.MaxInt64,
		"term_id": dto.TermID,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertState, stateArgs)
	if err != nil {
		d.log.Error("query execution failed", idAttr, slog.String("q", insertState))
		return err
	}
	d.log.Log(ds.Ctx, lf.LevelTrace, "entity insertion succeed", idAttr)
	return nil
}

func (d *daoPgx) UpdateType(source sd.Source, rec DefRec) error {
	ds := sd.MustConform[sd.SourcePgx](source)
	idAttr := slog.Any("defID", rec.DefID)
	d.log.Log(ds.Ctx, lf.LevelTrace, "entity update started", idAttr)
	dto, err := DataFromDefRec(rec)
	if err != nil {
		d.log.Error("model mapping failed", idAttr)
		return err
	}
	updateRoot := `
		update type_def_roots
		set def_rn = @def_rn,
			term_id = @term_id
		where def_id = @def_id
			and def_rn = @def_rn - 1`
	insertSnap := `
		insert into role_snaps (
			def_id, def_rn, title, term_id
		) values (
			@def_id, @def_rn, @title, @term_id
		)`
	args := pgx.NamedArgs{
		"def_id":  dto.DefID,
		"def_rn":  dto.DefRN,
		"title":   dto.Title,
		"term_id": dto.TermID,
	}
	ct, err := ds.Conn.Exec(ds.Ctx, updateRoot, args)
	if err != nil {
		d.log.Error("query execution failed", idAttr, slog.String("q", updateRoot))
		return err
	}
	if ct.RowsAffected() == 0 {
		d.log.Error("entity update failed", idAttr)
		return errOptimisticUpdate(rec.DefRN - 1)
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertSnap, args)
	if err != nil {
		d.log.Error("query execution failed", idAttr, slog.String("q", insertSnap))
		return err
	}
	d.log.Log(ds.Ctx, lf.LevelTrace, "entity update succeed", idAttr)
	return nil
}

func (d *daoPgx) SelectTypeRefs(source sd.Source) ([]DefRef, error) {
	ds := sd.MustConform[sd.SourcePgx](source)
	query := `
		SELECT
			def_id, def_rn, title
		FROM type_def_roots`
	rows, err := ds.Conn.Query(ds.Ctx, query)
	if err != nil {
		d.log.Error("query execution failed", slog.String("q", query))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[defRefDS])
	if err != nil {
		d.log.Error("rows collection failed")
		return nil, err
	}
	d.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return DataToDefRefs(dtos)
}

func (d *daoPgx) SelectTypeRecByID(source sd.Source, defID identity.ADT) (DefRec, error) {
	ds := sd.MustConform[sd.SourcePgx](source)
	idAttr := slog.Any("defID", defID)
	rows, err := ds.Conn.Query(ds.Ctx, selectById, defID.String())
	if err != nil {
		d.log.Error("query execution failed", idAttr, slog.String("q", selectById))
		return DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
	if err != nil {
		d.log.Error("row collection failed", idAttr)
		return DefRec{}, err
	}
	d.log.Log(ds.Ctx, lf.LevelTrace, "entity selection succeed", idAttr)
	return DataToDefRec(dto)
}

func (d *daoPgx) SelectTypeRecByQN(source sd.Source, defQN qualsym.ADT) (DefRec, error) {
	ds := sd.MustConform[sd.SourcePgx](source)
	fqnAttr := slog.Any("defQN", defQN)
	rows, err := ds.Conn.Query(ds.Ctx, selectByFQN, qualsym.ConvertToString(defQN))
	if err != nil {
		d.log.Error("query execution failed", fqnAttr, slog.String("q", selectByFQN))
		return DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
	if err != nil {
		d.log.Error("row collection failed", fqnAttr)
		return DefRec{}, err
	}
	d.log.Log(ds.Ctx, lf.LevelTrace, "entity selection succeed", fqnAttr)
	return DataToDefRec(dto)
}

func (d *daoPgx) SelectTypeRecsByIDs(source sd.Source, defIDs []identity.ADT) (_ []DefRec, err error) {
	ds := sd.MustConform[sd.SourcePgx](source)
	if len(defIDs) == 0 {
		return []DefRec{}, nil
	}
	query := `
		select
			def_id, def_rn, title, term_id, whole_id
		from type_def_roots
		where def_id = $1`
	batch := pgx.Batch{}
	for _, defID := range defIDs {
		if defID.IsEmpty() {
			return nil, identity.ErrEmpty
		}
		batch.Queue(query, defID.String())
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []defRecDS
	for _, defID := range defIDs {
		rows, err := br.Query()
		if err != nil {
			d.log.Error("query execution failed", slog.Any("defID", defID), slog.String("q", query))
		}
		defer rows.Close()
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
		if err != nil {
			d.log.Error("row collection failed", slog.Any("defID", defID))
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, err
	}
	d.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return DataToDefRecs(dtos)
}

func (d *daoPgx) SelectTypeEnv(source sd.Source, defQNs []qualsym.ADT) (map[qualsym.ADT]DefRec, error) {
	recs, err := d.SelectTypeRecsByQNs(source, defQNs)
	if err != nil {
		return nil, err
	}
	env := make(map[qualsym.ADT]DefRec, len(recs))
	for i, root := range recs {
		env[defQNs[i]] = root
	}
	return env, nil
}

func (d *daoPgx) SelectTypeRecsByQNs(source sd.Source, defQNs []qualsym.ADT) (_ []DefRec, err error) {
	ds := sd.MustConform[sd.SourcePgx](source)
	if len(defQNs) == 0 {
		return []DefRec{}, nil
	}
	batch := pgx.Batch{}
	for _, defQN := range defQNs {
		batch.Queue(selectByFQN, qualsym.ConvertToString(defQN))
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var dtos []defRecDS
	for _, defQN := range defQNs {
		rows, err := br.Query()
		if err != nil {
			d.log.Error("query execution failed", slog.Any("defQN", defQN), slog.String("q", selectByFQN))
		}
		defer rows.Close()
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDS])
		if err != nil {
			d.log.Error("row collection failed", slog.Any("defQN", defQN))
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, err
	}
	d.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("dtos", dtos))
	return DataToDefRecs(dtos)
}

func (d *daoPgx) InsertTerm(source sd.Source, rec TermRec) (err error) {
	ds := sd.MustConform[sd.SourcePgx](source)
	dto := dataFromTermRec(rec)
	query := `
		INSERT INTO type_def_states (
			id, kind, from_id, spec
		) VALUES (
			@id, @kind, @from_id, @spec
		)`
	batch := pgx.Batch{}
	for _, st := range dto.States {
		sa := pgx.NamedArgs{
			"id":      st.TermID,
			"kind":    st.K,
			"from_id": st.FromID,
			"spec":    st.Spec,
		}
		batch.Queue(query, sa)
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for range dto.States {
		_, err = br.Exec()
		if err != nil {
			d.log.Error("query execution failed", slog.Any("id", rec.Ident()), slog.String("q", query))
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (d *daoPgx) SelectTermRecByID(source sd.Source, defID identity.ADT) (TermRec, error) {
	ds := sd.MustConform[sd.SourcePgx](source)
	idAttr := slog.Any("defID", defID)
	query := `
		WITH RECURSIVE top_states AS (
			SELECT
				rs.*
			FROM type_def_states rs
			WHERE id = $1
			UNION ALL
			SELECT
				bs.*
			FROM type_def_states bs, top_states ts
			WHERE bs.from_id = ts.id
		)
		SELECT * FROM top_states`
	rows, err := ds.Conn.Query(ds.Ctx, query, defID.String())
	if err != nil {
		d.log.Error("query execution failed", idAttr, slog.String("q", query))
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[stateDS])
	if err != nil {
		d.log.Error("row collection failed", idAttr)
		return nil, err
	}
	if len(dtos) == 0 {
		d.log.Error("entity selection failed", idAttr)
		return nil, fmt.Errorf("no rows selected")
	}
	d.log.Log(ds.Ctx, lf.LevelTrace, "entity selection succeed", slog.Any("dtos", dtos))
	type_def_states := make(map[string]stateDS, len(dtos))
	for _, dto := range dtos {
		type_def_states[dto.TermID] = dto
	}
	return statesToTermRec(type_def_states, type_def_states[defID.String()])
}

func (d *daoPgx) SelectTermEnv(source sd.Source, recIDs []identity.ADT) (map[identity.ADT]TermRec, error) {
	recs, err := d.SelectTermRecsByIDs(source, recIDs)
	if err != nil {
		return nil, err
	}
	env := make(map[identity.ADT]TermRec, len(recs))
	for _, rec := range recs {
		env[rec.Ident()] = rec
	}
	return env, nil
}

func (d *daoPgx) SelectTermRecsByIDs(source sd.Source, recIDs []identity.ADT) (_ []TermRec, err error) {
	ds := sd.MustConform[sd.SourcePgx](source)
	batch := pgx.Batch{}
	for _, rid := range recIDs {
		batch.Queue(selectByID, rid.String())
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	var recs []TermRec
	for _, recID := range recIDs {
		idAttr := slog.Any("recID", recID)
		rows, err := br.Query()
		if err != nil {
			d.log.Error("query execution failed", idAttr, slog.String("q", selectByID))
		}
		defer rows.Close()
		dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[stateDS])
		if err != nil {
			d.log.Error("rows collection failed", idAttr)
		}
		if len(dtos) == 0 {
			d.log.Error("entity selection failed", idAttr)
			return nil, ErrDoesNotExist(recID)
		}
		rec, err := dataToTermRec(&termRecDS{recID.String(), dtos})
		if err != nil {
			d.log.Error("model mapping failed", idAttr)
			return nil, err
		}
		recs = append(recs, rec)
	}
	d.log.Log(ds.Ctx, lf.LevelTrace, "entities selection succeed", slog.Any("recs", recs))
	return recs, err
}

const (
	selectByFQN = `
		select
			rr.def_id,
			rr.def_rn,
			rr.title,
			rs.term_id,
			null as whole_id
		from type_def_roots rr
		left join aliases a
			on a.id = rr.def_id
			and a.from_rn >= rr.def_rn
			and a.to_rn > rr.def_rn
		left join type_def_states rs
			on rs.def_id = rr.def_id
			and rs.from_rn >= rr.def_rn
			and rs.to_rn > rr.def_rn
		where a.sym = $1`

	selectById = `
		select
			rr.def_id,
			rr.def_rn,
			rr.title,
			rs.term_id,
			null as whole_id
		from type_def_roots rr
		left join type_def_states rs
			on rs.def_id = rr.def_id
			and rs.from_rn >= rr.def_rn
			and rs.to_rn > rr.def_rn
		where rr.def_id = $1`

	selectByID = `
		WITH RECURSIVE state_tree AS (
			SELECT root.*
			FROM type_def_states root
			WHERE id = $1
			UNION ALL
			SELECT child.*
			FROM type_def_states child, state_tree parent
			WHERE child.from_id = parent.id
		)
		SELECT * FROM state_tree`
)
