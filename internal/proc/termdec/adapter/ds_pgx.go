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

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termsem"

	"orglang/go-engine/internal/proc/ent"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

type decRootRef struct {
	DescID string `db:"desc_id"`
	Rev    int64  `db:"rev"`
	Title  string `db:"title"`
}

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) termdec.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

// for compilation purposes
func newRepo() termdec.Repo {
	return new(daoEnt)
}

func (dao *daoEnt) AddRec(ctx context.Context, conn tx.Conn, rec termdec.DecRec) error {
	refAttr := slog.Any("ref", rec.TermRef)
	dto, err := DataFromDecRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}
	_, err = dao.client.TermDec.Create().
		SetID(dto.TermID).
		SetTermRn(dto.TermRN).
		SetLiabVar(dto.LiabVar).
		SetAssetVars(dto.AssetVars).
		Save(entpgx.WithConn(ctx, conn.Pgx))
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}
	return nil
}

func (dao *daoPgx) GetSnap(ctx context.Context, conn tx.Conn, ref termsem.SemRef) (termdec.DecSnap, error) {
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
		return termdec.DecSnap{}, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entitiy selection succeed", slog.Any("dto", dto))
	return DataToDecSnap(dto)
}

func (dao *daoPgx) SelectEnv(ctx context.Context, conn tx.Conn, ids []identity.ADT) (map[identity.ADT]termdec.DecRec, error) {
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

func (dao *daoPgx) GetRecs(ctx context.Context, conn tx.Conn, ids []identity.ADT) (_ []termdec.DecRec, err error) {
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
	var dtos []DecRec
	for _, rid := range ids {
		rows, err := br.Query()
		if err != nil {
			dao.log.Error("query execution failed", slog.Any("id", rid), slog.String("q", selectByRef))
		}
		dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DecRec])
		if err != nil {
			dao.log.Error("row scanning failed", slog.Any("id", rid))
		}
		dtos = append(dtos, dto)
	}
	if err != nil {
		return nil, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "selection succeed", slog.Any("dtos", dtos))
	return DataToDecRecs(dtos)
}

func (dao *daoPgx) GetRefs(ctx context.Context, conn tx.Conn) ([]termsem.SemRef, error) {
	query := `
		select
			desc_id, rev, title
		from dec_roots`
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
		id, convErr := identity.ConvertFromString(dto.DescID)
		if convErr != nil {
			return nil, convErr
		}
		refs[i] = termsem.SemRef{TermID: id, TermRN: seqnum.ConvertFromInt(dto.Rev)}
	}
	return refs, nil
}

const (
	selectByRef = `
		select
			pd.desc_id,
			ds.desc_rn,
			pd.liab_var,
			pd.asset_vars
		from proc_term_decs pd
		left join desc_sems ds
			on ds.desc_id = pd.desc_id
		where pd.desc_id = $1`
)
