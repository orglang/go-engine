package adapter

import (
	"context"
	"database/sql"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	commactadapter "orglang/go-engine/internal/pool/commact/adapter"
	commconn "orglang/go-engine/internal/pool/commconn/core"
)

type daoPgx struct {
	qb  commconn.QueryBuilder
	log *slog.Logger
}

func NewDaoPgx(qb commconn.QueryBuilder, log *slog.Logger) commconn.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{qb, log.With(name)}
}

// for compilation purposes
func newRepo() commconn.Repo { return new(daoPgx) }

func (dao *daoPgx) AddRec(ctx context.Context, conn tx.Conn, rec commconn.ConnRec) error {
	dto := DataFromRec(rec)
	refAttr := slog.Any("conn_id", rec.ConnID)
	querySQL, args := dao.qb.InsertRec(dto)
	_, err := conn.Pgx.Exec(ctx, querySQL, args...)
	if err != nil {
		dao.log.Error("Connection insertion failed", refAttr, slog.String("sql", querySQL))
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection insertion succeeded", slog.Any("dto", dto))
	return nil
}

func (dao *daoPgx) ModifyRec(ctx context.Context, conn tx.Conn, mod commconn.ConnMod) error {
	refAttr := slog.Any("conn_id", mod.ConnID)
	if (mod.OffsetNr != nil && !mod.OffsetNr.IsEmpty()) || len(mod.Acts) > 0 {
		dto := DataFromMod(mod)
		if len(mod.Acts) > 0 {
			// The provider CAS advances comp_rn exactly once for this operation.
			// Mirror that revision as the Act number; no independent sequence
			// update is needed (Take appends at most one Act per operation).
			dto.ActNr = sql.Null[int64]{V: int64(mod.ProviderCompRN) + 1, Valid: true}
		}
		querySQL, args := dao.qb.UpdateRec(dto)
		tag, err := conn.Pgx.Exec(ctx, querySQL, args...)
		if err != nil {
			dao.log.Error("Connection offset update failed", refAttr, slog.String("sql", querySQL))
			return err
		}
		if tag.RowsAffected() == 0 {
			return commconn.ErrMissingConn(mod.ConnID)
		}
	}
	for _, act := range mod.Acts {
		dto := commactadapter.DataFromStepRec(act)
		dto.ActNr = int64(mod.ProviderCompRN) + 1
		querySQL, args := dao.qb.InsertActRec(dto)
		if _, err := conn.Pgx.Exec(ctx, querySQL, args...); err != nil {
			dao.log.Error("Act insertion failed", refAttr, slog.String("sql", querySQL))
			return err
		}
	}
	dao.log.Log(ctx, lf.LevelTrace, "Connection modification succeeded", refAttr)
	return nil
}

func (dao *daoPgx) GetSnapByQry(ctx context.Context, conn tx.Conn, qry commconn.ConnQry) (commconn.ConnSnap, error) {
	refAttr := slog.Any("conn_id", qry.ConnID)
	qryDTO := DataFromQry(qry)
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot selection started", slog.Any("qry", qryDTO))
	querySQL, args := dao.qb.SelectSnap(qryDTO)
	rows, err := conn.Pgx.Query(ctx, querySQL, args...)
	if err != nil {
		dao.log.Error("Connection snapshot query failed", refAttr, slog.String("sql", querySQL))
		return commconn.ConnSnap{}, err
	}
	defer rows.Close()
	snapDTO, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[commconn.ConnSnapDS])
	if err != nil {
		dao.log.Error("Connection snapshot scan failed", refAttr)
		return commconn.ConnSnap{}, err
	}
	snap, err := DataToSnap(snapDTO)
	if err != nil {
		dao.log.Error("Connection snapshot conversion failed", refAttr)
		return commconn.ConnSnap{}, err
	}
	providerCompID, err := identity.ConvertFromString(snapDTO.ProviderCompID)
	if err != nil {
		return commconn.ConnSnap{}, err
	}
	snap.ProviderCompID = providerCompID
	snap.ProviderCompRN = seqnum.ADT(snapDTO.ProviderCompRN)
	dao.log.Log(ctx, lf.LevelTrace, "Connection snapshot loaded", slog.Any("dto", snapDTO))
	return snap, nil
}
