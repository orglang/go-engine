package adapter

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/progroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termvar"
	termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/proc/ent"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

type decRootRef struct {
	TermID string `db:"term_id"`
}

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) termdec.Repo {
	return &daoEnt{client: client, log: log}
}

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec termdec.DecRec) error {
	dto, err := dataFromTermDecRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", slog.Any("termID", rec.TermID))
		return err
	}
	_, err = dao.client.TermDec.Create().SetID(dto.TermID).SetProgID(dto.ProgID).SetLiabVar(dto.LiabVar).SetAssetVars(dto.AssetVars).Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", slog.Any("termID", rec.TermID))
		return err
	}
	const bindQuery = `
INSERT INTO proc_prog_binds (prog_qn, prog_id, kind)
VALUES ($1, $2, 2)
`
	if _, err := conn.Pgx.Exec(ctx, bindQuery, uniqsym.ConvertToString(rec.TermQN), rec.TermID.String()); err != nil {
		dao.log.Error("program binding insertion failed", slog.Any("termID", rec.TermID), slog.Any("error", err))
		return err
	}
	return nil
}

func (dao *daoEnt) GetSnap(ctx context.Context, conn tx.Conn, id identity.ADT) (termdec.DecSnap, error) {
	rows, err := conn.Pgx.Query(ctx, selectByID, id.String())
	if err != nil {
		dao.log.Error("query execution failed", slog.Any("termID", id))
		return termdec.DecSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DecSnap])
	if err != nil {
		dao.log.Error("row scanning failed", slog.Any("termID", id))
		return termdec.DecSnap{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entity selection succeed", slog.Any("dto", dto))
	return dataToTermDecSnap(dto)
}

func (dao *daoEnt) SelectEnv(ctx context.Context, conn tx.Conn, ids []identity.ADT) (map[identity.ADT]termdec.DecRec, error) {
	recs, err := dao.GetRecs(ctx, conn, ids)
	if err != nil {
		return nil, err
	}
	env := make(map[identity.ADT]termdec.DecRec, len(recs))
	for _, rec := range recs {
		env[rec.TermID] = rec
	}
	return env, nil
}

func (dao *daoEnt) GetRecs(ctx context.Context, conn tx.Conn, ids []identity.ADT) (_ []termdec.DecRec, err error) {
	if len(ids) == 0 {
		return []termdec.DecRec{}, nil
	}
	batch := pgx.Batch{}
	for _, id := range ids {
		if id.IsEmpty() {
			return nil, identity.ErrEmpty
		}
		batch.Queue(selectByID, id.String())
	}
	results := conn.Pgx.SendBatch(ctx, &batch)
	defer func() { err = errors.Join(err, results.Close()) }()
	dtos := make([]DecRec, 0, len(ids))
	for _, id := range ids {
		rows, queryErr := results.Query()
		if queryErr != nil {
			dao.log.Error("query execution failed", slog.Any("id", id))
			return nil, queryErr
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DecRec])
		if scanErr != nil {
			dao.log.Error("row scanning failed", slog.Any("id", id))
			return nil, scanErr
		}
		dtos = append(dtos, dto)
	}
	recs := make([]termdec.DecRec, len(dtos))
	for i, dto := range dtos {
		recs[i], err = dataToTermDecRec(dto)
		if err != nil {
			return nil, err
		}
	}
	return recs, nil
}

func (dao *daoEnt) GetQNByID(ctx context.Context, conn tx.Conn, termID identity.ADT) (uniqsym.ADT, error) {
	const query = `
SELECT prog_qn
FROM proc_prog_binds
WHERE prog_id = $1
`
	rows, err := conn.Pgx.Query(ctx, query, termID.String())
	if err != nil {
		return uniqsym.ADT{}, err
	}
	defer rows.Close()
	qn, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[string])
	if err != nil {
		return uniqsym.ADT{}, err
	}
	return uniqsym.ConvertFromString(qn)
}

func (dao *daoEnt) GetIDs(ctx context.Context, conn tx.Conn) ([]identity.ADT, error) {
	const query = `
SELECT term_id
FROM proc_prog_terms
`
	rows, err := conn.Pgx.Query(ctx, query)
	if err != nil {
		dao.log.Error("query execution failed")
		return nil, err
	}
	defer rows.Close()
	values, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		dao.log.Error("rows scanning failed")
		return nil, err
	}
	ids := make([]identity.ADT, len(values))
	for i, value := range values {
		ids[i], err = identity.ConvertFromString(value)
		if err != nil {
			return nil, err
		}
	}
	return ids, nil
}

const selectByID = `
SELECT
    pd.term_id,
    pd.prog_id,
    dr.prog_rn,
    pd.liab_var,
    pd.asset_vars
FROM proc_prog_terms AS pd
JOIN proc_prog_roots AS dr
    ON dr.prog_id = pd.prog_id
WHERE pd.term_id = $1
`

func dataFromTermDecRec(rec termdec.DecRec) (DecRec, error) {
	dto := DecRec{TermID: identity.ConvertToString(rec.TermID), ProgID: identity.ConvertToString(rec.ProgRef.ProgID), ProgRN: seqnum.ConvertToInt(rec.ProgRef.ProgRN)}
	var err error
	dto.LiabVar, err = termvaradapter.DataFromVarRec(rec.LiabVar)
	if err != nil {
		return DecRec{}, err
	}
	if rec.AssetVars != nil {
		dto.AssetVars = make([]termvaradapter.VarRec, len(rec.AssetVars))
		for i, variable := range rec.AssetVars {
			dto.AssetVars[i], err = termvaradapter.DataFromVarRec(variable)
			if err != nil {
				return DecRec{}, err
			}
		}
	}
	return dto, nil
}

func dataToTermDecRec(dto DecRec) (termdec.DecRec, error) {
	termID, err := identity.ConvertFromString(dto.TermID)
	if err != nil {
		return termdec.DecRec{}, err
	}
	descID, err := identity.ConvertFromString(dto.ProgID)
	if err != nil {
		return termdec.DecRec{}, err
	}
	rec := termdec.DecRec{TermID: termID, ProgRef: progroot.RootRef{ProgID: descID, ProgRN: seqnum.ConvertFromInt(dto.ProgRN)}}
	rec.LiabVar, err = termvaradapter.DataToVarRec(dto.LiabVar)
	if err != nil {
		return termdec.DecRec{}, err
	}
	if dto.AssetVars != nil {
		rec.AssetVars = make([]termvar.VarRec, len(dto.AssetVars))
		for i, variable := range dto.AssetVars {
			rec.AssetVars[i], err = termvaradapter.DataToVarRec(variable)
			if err != nil {
				return termdec.DecRec{}, err
			}
		}
	}
	return rec, nil
}

func dataToTermDecSnap(dto DecSnap) (termdec.DecSnap, error) {
	termID, err := identity.ConvertFromString(dto.TermID)
	if err != nil {
		return termdec.DecSnap{}, err
	}
	descID, err := identity.ConvertFromString(dto.ProgID)
	if err != nil {
		return termdec.DecSnap{}, err
	}
	snap := termdec.DecSnap{TermID: termID, ProgRef: progroot.RootRef{ProgID: descID, ProgRN: seqnum.ConvertFromInt(dto.ProgRN)}}
	snap.LiabVar, err = termvaradapter.DataToVarRec(dto.LiabVar)
	if err != nil {
		return termdec.DecSnap{}, err
	}
	if dto.AssetVars != nil {
		snap.AssetVars = make([]termvar.VarRec, len(dto.AssetVars))
		for i, variable := range dto.AssetVars {
			snap.AssetVars[i], err = termvaradapter.DataToVarRec(variable)
			if err != nil {
				return termdec.DecSnap{}, err
			}
		}
	}
	return snap, nil
}
