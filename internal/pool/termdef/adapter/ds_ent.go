package adapter

import (
	"context"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) termdef.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() termdef.Repo {
	return new(daoEnt)
}

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec termdef.DefRec) error {
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

func (dao *daoEnt) GetRecByQN(ctx context.Context, conn tx.Conn, qn uniqsym.ADT) (termdef.DefRec, error) {
	qnAttr := slog.Any("qn", qn)
	const query = `
SELECT def.term_id, def.desc_id, root.desc_rn, def.liab_var, def.asset_vars
FROM pool_term_decs AS def
JOIN pool_desc_roots AS root ON root.desc_id = def.desc_id
JOIN pool_desc_binds AS bind ON bind.desc_id = def.term_id
WHERE bind.desc_qn = $1
`
	rows, err := conn.Pgx.Query(ctx, query, uniqsym.ConvertToString(qn))
	if err != nil {
		dao.log.Error("query execution failed", qnAttr)
		return termdef.DefRec{}, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DefRec])
	if err != nil {
		dao.log.Error("rows scanning failed", qnAttr)
		return termdef.DefRec{}, err
	}
	rec, err := DataToDecRec(dto)
	if err != nil {
		dao.log.Error("model conversion failed", qnAttr)
		return termdef.DefRec{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "getting succeed", qnAttr)
	return rec, nil
}
