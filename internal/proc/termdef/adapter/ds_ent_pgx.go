package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/uniqsym"

	"orglang/go-engine/internal/proc/ent"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	termdef "orglang/go-engine/internal/proc/termdef/core"
)

type decRootRef struct {
	TermID string `db:"term_id"`
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
}

type defRecDTO struct {
	TermID string `db:"term_id"`
	ExpCK  int64  `db:"exp_ck"`
}

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) *daoEnt {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() termdef.Repo {
	return new(daoEnt)
}

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec termdec.DecRec) error {
	refAttr := slog.Any("ref", rec.TermRef)
	dto, err := DataFromTermDecRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}
	_, err = dao.client.TermDec.Create().
		SetID(dto.TermID).
		SetDescID(dto.DescID).
		SetLiabVar(dto.LiabVar).
		SetAssetVars(dto.AssetVars).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}
	return nil
}

func (dao *daoEnt) SetExpCK(ctx context.Context, conn tx.Conn, termID identity.ADT, expCK contkey.ADT) error {
	const query = `
UPDATE proc_desc_terms
SET exp_ck = $1
WHERE term_id = $2
`
	result, err := conn.Pgx.Exec(ctx, query, int64(expCK), termID.String())
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return termdef.ErrDoesNotExist(termID)
	}
	return nil
}

func (dao *daoEnt) GetSnap(ctx context.Context, conn tx.Conn, ref termsem.SemRef) (termdec.DecSnap, error) {
	refAttr := slog.Any("ref", ref)
	rows, err := conn.Pgx.Query(ctx, selectByRef, ref.TermID.String())
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return termdec.DecSnap{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DecSnap])
	if err != nil {
		dao.log.Error("row scanning failed", refAttr)
		return termdef.DecSnap{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entity selection succeed", slog.Any("dto", dto))
	return DataToTermDecSnap(dto)
}

func (dao *daoEnt) SelectEnv(ctx context.Context, conn tx.Conn, ids []identity.ADT) (map[identity.ADT]termdec.DecRec, error) {
	decs, err := dao.GetRecs(ctx, conn, ids)
	if err != nil {
		return nil, err
	}
	env := make(map[identity.ADT]termdec.DecRec, len(decs))
	for _, dec := range decs {
		env[dec.TermRef.TermID] = dec
	}
	return env, nil
}

func (dao *daoEnt) GetRecs(ctx context.Context, conn tx.Conn, ids []identity.ADT) (_ []termdec.DecRec, err error) {
	if len(ids) == 0 {
		return []termdec.DecRec{}, nil
	}
	batch := pgx.Batch{}
	for _, rid := range ids {
		if rid.IsEmpty() {
			return nil, identity.ErrEmpty
		}
		batch.Queue(selectByRef, rid.String())
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()

	dtos := make([]DecRec, 0, len(ids))
	for _, rid := range ids {
		idAttr := slog.Any("id", rid)
		rows, queryErr := br.Query()
		if queryErr != nil {
			dao.log.Error("query execution failed", idAttr, slog.String("q", selectByRef))
			return nil, queryErr
		}
		dto, scanErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DecRec])
		if scanErr != nil {
			dao.log.Error("row scanning failed", idAttr)
			return nil, scanErr
		}
		dtos = append(dtos, dto)
	}
	dao.log.Log(ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	return DataToTermDecRecs(dtos)
}

func (dao *daoEnt) GetRecByID(ctx context.Context, conn tx.Conn, termID identity.ADT) (termdef.DefRec, error) {
	if termID.IsEmpty() {
		return termdef.DefRec{}, identity.ErrEmpty
	}
	const query = `
SELECT pd.term_id, COALESCE(pd.exp_ck, 0) AS exp_ck
FROM proc_desc_terms AS pd
WHERE pd.term_id = $1
`
	rows, err := conn.Pgx.Query(ctx, query, termID.String())
	if err != nil {
		return termdef.DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[defRecDTO])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return termdef.DefRec{}, termdef.ErrTermMissing(termID)
		}
		return termdef.DefRec{}, err
	}
	expCK, err := contkey.ConvertFromInt(dto.ExpCK)
	if err != nil {
		return termdef.DefRec{}, err
	}
	return termdef.DefRec{TermID: termID, ExpCK: expCK}, nil
}

func (dao *daoEnt) GetQNByID(ctx context.Context, conn tx.Conn, termID identity.ADT) (uniqsym.ADT, error) {
	const query = `
SELECT desc_qn
FROM proc_desc_binds
WHERE desc_id = $1
`
	var qn string
	if err := conn.Pgx.QueryRow(ctx, query, termID.String()).Scan(&qn); err != nil {
		return "", err
	}
	return uniqsym.ConvertFromString(qn)
}

func (dao *daoEnt) GetRefs(ctx context.Context, conn tx.Conn) ([]termsem.SemRef, error) {
	query := `
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
    COALESCE(pd.exp_ck, 0) AS exp_ck,
    pd.liab_var,
    pd.asset_vars
FROM proc_desc_terms AS pd
JOIN proc_desc_roots AS dr
    ON dr.desc_id = pd.desc_id
WHERE pd.term_id = $1
`
