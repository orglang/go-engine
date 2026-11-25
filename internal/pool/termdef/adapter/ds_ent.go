package adapter

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

type defRecDTO struct {
	TermID string `db:"term_id"`
	DescID string `db:"desc_id"`
	DescRN int64  `db:"desc_rn"`
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

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec termdef.DecRec) error {
	refAttr := slog.Any("ref", rec.TermRef)
	dao.log.Log(ctx, lf.LevelTrace, "addition started", refAttr)

	dto, err := DataFromDecRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}

	_, err = dao.client.TermDef.Create().
		SetID(dto.TermID).
		SetDescID(dto.DescID).
		SetLiabVar(dto.LiabVar).
		SetAssetVars(dto.AssetVars).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}

	dao.log.Log(ctx, lf.LevelTrace, "addition succeed", refAttr)
	return nil
}

func (dao *daoEnt) SetExpCK(ctx context.Context, conn tx.Conn, termID identity.ADT, expCK contkey.ADT) error {
	const query = `
UPDATE pool_desc_terms
SET exp_ck = $1
WHERE term_id = $2
`
	result, err := conn.Pgx.Exec(ctx, query, int64(expCK), termID.String())
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("pool term %v does not exist", termID)
	}
	return nil
}

func (dao *daoEnt) GetRecByID(ctx context.Context, conn tx.Conn, termID identity.ADT) (termdef.DefRec, error) {
	const query = `
SELECT
    def.term_id,
    def.desc_id,
    root.desc_rn,
    COALESCE(def.exp_ck, 0) AS exp_ck
FROM pool_desc_terms AS def
JOIN pool_desc_roots AS root
    ON root.desc_id = def.desc_id
WHERE def.term_id = $1
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
	descID, err := identity.ConvertFromString(dto.DescID)
	if err != nil {
		return termdef.DefRec{}, err
	}
	expCK, err := contkey.ConvertFromInt(dto.ExpCK)
	if err != nil {
		return termdef.DefRec{}, err
	}
	return termdef.DefRec{TermRef: termsem.SemRef{TermID: termID, DescID: descID, DescRN: seqnum.ConvertFromInt(dto.DescRN)}, ExpCK: expCK}, nil
}

func (dao *daoEnt) GetRecByQN(ctx context.Context, conn tx.Conn, qn uniqsym.ADT) (termdef.DecRec, error) {
	qnAttr := slog.Any("qn", qn)
	const query = `
SELECT def.term_id, def.desc_id, root.desc_rn, COALESCE(def.exp_ck, 0) AS exp_ck, def.liab_var, def.asset_vars
FROM pool_desc_terms AS def
JOIN pool_desc_roots AS root ON root.desc_id = def.desc_id
JOIN pool_desc_binds AS bind ON bind.desc_id = def.term_id
WHERE bind.desc_qn = $1
`
	rows, err := conn.Pgx.Query(ctx, query, uniqsym.ConvertToString(qn))
	if err != nil {
		dao.log.Error("query execution failed", qnAttr)
		return termdef.DecRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DecRec])
	if err != nil {
		dao.log.Error("rows scanning failed", qnAttr)
		return termdef.DecRec{}, err
	}
	rec, err := DataToDecRec(dto)
	if err != nil {
		dao.log.Error("model conversion failed", qnAttr)
		return termdef.DecRec{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", qnAttr)
	return rec, nil
}
