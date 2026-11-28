package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termvar"
	termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	termdec "orglang/go-engine/internal/pool/termdec/core"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) termdec.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec termdec.DecRec) error {
	refAttr := slog.Any("termID", rec.TermID)
	liabVar, err := termvaradapter.DataFromVarRec(rec.LiabVar)
	if err != nil {
		return err
	}
	assetVars := make([]termvaradapter.VarRec, len(rec.AssetVars))
	for i, variable := range rec.AssetVars {
		assetVars[i], err = termvaradapter.DataFromVarRec(variable)
		if err != nil {
			return err
		}
	}
	_, err = dao.client.TermDec.Create().
		SetID(rec.TermID.String()).
		SetDescID(rec.DescRef.DescID.String()).
		SetLiabVar(liabVar).
		SetAssetVars(assetVars).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}
	const bindQuery = `
INSERT INTO pool_desc_binds (desc_qn, desc_id, kind)
VALUES ($1, $2, 2)
`
	if _, err := conn.Pgx.Exec(ctx, bindQuery, uniqsym.ConvertToString(rec.TermQN), rec.TermID.String()); err != nil {
		dao.log.Error("description binding insertion failed", refAttr, slog.Any("error", err))
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "declaration addition succeeded", refAttr)
	return nil
}

func (dao *daoEnt) GetRecByQN(ctx context.Context, conn tx.Conn, qn uniqsym.ADT) (termdec.DecRec, error) {
	const query = `
SELECT def.term_id, def.desc_id, root.desc_rn, def.liab_var, def.asset_vars
FROM pool_desc_terms AS def
JOIN pool_desc_roots AS root ON root.desc_id = def.desc_id
JOIN pool_desc_binds AS bind ON bind.desc_id = def.term_id
WHERE bind.desc_qn = $1
`
	rows, err := conn.Pgx.Query(ctx, query, uniqsym.ConvertToString(qn))
	if err != nil {
		return termdec.DecRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DecRec])
	if err != nil {
		return termdec.DecRec{}, err
	}
	termID, err := identity.ConvertFromString(dto.TermID)
	if err != nil {
		return termdec.DecRec{}, err
	}
	descID, err := identity.ConvertFromString(dto.DescID)
	if err != nil {
		return termdec.DecRec{}, err
	}
	liabVar, err := termvaradapter.DataToVarRec(dto.LiabVar)
	if err != nil {
		return termdec.DecRec{}, err
	}
	assetVars := make([]termvar.VarRec, len(dto.AssetVars))
	for i, variable := range dto.AssetVars {
		assetVars[i], err = termvaradapter.DataToVarRec(variable)
		if err != nil {
			return termdec.DecRec{}, err
		}
	}
	return termdec.DecRec{
		TermQN:    qn,
		TermID:    termID,
		DescRef:   descroot.RootRef{DescID: descID, DescRN: seqnum.ConvertFromInt(dto.DescRN)},
		LiabVar:   liabVar,
		AssetVars: assetVars,
	}, nil
}

func (dao *daoEnt) GetQNByID(ctx context.Context, conn tx.Conn, termID identity.ADT) (uniqsym.ADT, error) {
	const query = `
SELECT desc_qn
FROM pool_desc_binds
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

var _ termdec.Repo = (*daoEnt)(nil)
