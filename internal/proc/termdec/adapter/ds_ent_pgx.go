package adapter

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/termvar"
	termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/proc/ent"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

type decRootRef struct {
	TermID string `db:"term_id"`
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
}

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) *daoEnt {
	return &daoEnt{client: client, log: log}
}

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec termdec.DecRec) error {
	dto, err := dataFromTermDecRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", slog.Any("ref", rec.TermRef))
		return err
	}
	_, err = dao.client.TermDec.Create().
		SetID(dto.TermID).
		SetDescID(dto.DescID).
		SetLiabVar(dto.LiabVar).
		SetAssetVars(dto.AssetVars).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", slog.Any("ref", rec.TermRef))
		return err
	}
	return nil
}

func (dao *daoEnt) GetSnap(ctx context.Context, conn tx.Conn, ref termsem.SemRef) (termdec.DecSnap, error) {
	rows, err := conn.Pgx.Query(ctx, selectByRef, ref.TermID.String())
	if err != nil {
		dao.log.Error("query execution failed", slog.Any("ref", ref))
		return termdec.DecSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DecSnap])
	if err != nil {
		dao.log.Error("row scanning failed", slog.Any("ref", ref))
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
		env[rec.TermRef.TermID] = rec
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
		batch.Queue(selectByRef, id.String())
	}
	results := conn.Pgx.SendBatch(ctx, &batch)
	defer func() { err = errors.Join(err, results.Close()) }()

	dtos := make([]DecRec, 0, len(ids))
	for _, id := range ids {
		rows, queryErr := results.Query()
		if queryErr != nil {
			dao.log.Error("query execution failed", slog.Any("id", id), slog.String("q", selectByRef))
			return nil, queryErr
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DecRec])
		if scanErr != nil {
			dao.log.Error("row scanning failed", slog.Any("id", id))
			return nil, scanErr
		}
		dtos = append(dtos, dto)
	}
	dao.log.Log(ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
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
SELECT desc_qn
FROM proc_desc_binds
WHERE desc_id = $1
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

func (dao *daoEnt) GetRefs(ctx context.Context, conn tx.Conn) ([]termsem.SemRef, error) {
	const query = `
SELECT td.term_id, td.desc_id, dr.desc_rn
FROM proc_desc_terms AS td
JOIN proc_desc_roots AS dr ON dr.desc_id = td.desc_id
`
	rows, err := conn.Pgx.Query(ctx, query)
	if err != nil {
		dao.log.Error("query execution failed")
		return nil, err
	}
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[decRootRef])
	if err != nil {
		dao.log.Error("rows scanning failed")
		return nil, err
	}
	refs := make([]termsem.SemRef, len(dtos))
	for i, dto := range dtos {
		id, convErr := identity.ConvertFromString(dto.TermID)
		if convErr != nil {
			return nil, convErr
		}
		descID, convErr := identity.ConvertFromString(dto.DescID)
		if convErr != nil {
			return nil, convErr
		}
		refs[i] = termsem.SemRef{TermID: id, DescID: descID, DescRN: seqnum.ConvertFromInt(dto.DescRN)}
	}
	return refs, nil
}

const selectByRef = `
SELECT
    pd.term_id,
    pd.desc_id,
    dr.desc_rn,
    pd.liab_var,
    pd.asset_vars
FROM proc_desc_terms AS pd
JOIN proc_desc_roots AS dr
    ON dr.desc_id = pd.desc_id
WHERE pd.term_id = $1
`

func dataFromTermDecRec(rec termdec.DecRec) (DecRec, error) {
	dto := DecRec{
		TermID: identity.ConvertToString(rec.TermRef.TermID),
		DescID: identity.ConvertToString(rec.TermRef.DescID),
		DescRN: seqnum.ConvertToInt(rec.TermRef.DescRN),
	}
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
	ref, err := termsemRef(dto.TermID, dto.DescID, dto.DescRN)
	if err != nil {
		return termdec.DecRec{}, err
	}
	rec := termdec.DecRec{TermRef: ref}
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
	ref, err := termsemRef(dto.TermID, dto.DescID, dto.DescRN)
	if err != nil {
		return termdec.DecSnap{}, err
	}
	snap := termdec.DecSnap{TermRef: ref}
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

func termsemRef(termID, descID string, descRN int64) (termsem.SemRef, error) {
	id, err := identity.ConvertFromString(termID)
	if err != nil {
		return termsem.SemRef{}, err
	}
	rootID, err := identity.ConvertFromString(descID)
	if err != nil {
		return termsem.SemRef{}, err
	}
	return termsem.SemRef{TermID: id, DescID: rootID, DescRN: seqnum.ConvertFromInt(descRN)}, nil
}
