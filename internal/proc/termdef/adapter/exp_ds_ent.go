package adapter

import (
	"context"
	"encoding/json"
	"encoding/json/jsontext"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	termexpent "orglang/go-engine/internal/proc/ent/termexp"
	termdef "orglang/go-engine/internal/proc/termdef/core"
)

func (dao *daoEnt) InsertSpec(ctx context.Context, conn tx.Conn, spec termdef.ExpSpec) error {
	rec, err := termdef.ConvertSpecToRec(spec)
	if err != nil {
		return err
	}
	dto, err := dataFromExpRec(rec)
	if err != nil {
		return err
	}

	builders := make([]*ent.TermExpCreate, 0, len(dto.States))
	for _, state := range dto.States {
		raw, err := json.Marshal(state.Spec)
		if err != nil {
			return err
		}
		builders = append(builders, dao.client.TermExp.Create().
			SetID(state.ExpCK).
			SetSupExpCk(state.SupExpCK).
			SetKind(state.K).
			SetSpec(jsontext.Value(raw)))
	}

	txCtx := entpgx.WithConn(ctx, conn.Pgx)
	err = dao.client.TermExp.CreateBulk(builders...).
		OnConflictColumns(termexpent.FieldID).
		DoNothing().
		Exec(txCtx)
	if err != nil {
		dao.log.Error("process term Definition insertion failed", slog.Any("error", err))
		return err
	}
	dao.log.Log(ctx, lf.LevelTrace, "term Definition insertion succeeded", slog.Int("states", len(dto.States)))
	return nil
}

func (dao *daoEnt) SelectSpecByCK(ctx context.Context, conn tx.Conn, key contkey.ADT) (termdef.ExpSpec, error) {
	rows, err := conn.Pgx.Query(ctx, selectTermExpTree, int64(key))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	states, err := pgx.CollectRows(rows, pgx.RowToStructByName[ExpState])
	if err != nil {
		return nil, err
	}
	if len(states) == 0 {
		return nil, fmt.Errorf("process term Definition %d does not exist", key)
	}
	rec, err := dataToExpRec(ExpData{ExpCK: int64(key), States: states})
	if err != nil {
		return nil, err
	}
	return termdef.ConvertRecToSpec(rec), nil
}

const selectTermExpTree = `
WITH RECURSIVE exp_tree AS (
    SELECT top.*
    FROM proc_term_exps AS top
    WHERE top.exp_ck = $1
    UNION
    SELECT sub.*
    FROM proc_term_exps AS sub
    JOIN exp_tree AS sup ON sub.sup_exp_ck = sup.exp_ck
)
SELECT exp_ck, sup_exp_ck, kind, spec
FROM exp_tree
`
