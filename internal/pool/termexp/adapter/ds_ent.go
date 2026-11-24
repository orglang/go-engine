package adapter

import (
	"context"
	"encoding/json"
	"encoding/json/jsontext"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/lib/entpgx"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	termexpent "orglang/go-engine/internal/pool/ent/termexp"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

type daoEnt struct {
	client *ent.Client
	log    *slog.Logger
}

func NewDaoEnt(client *ent.Client, log *slog.Logger) termexp.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

func newRepo() termexp.Repo { return new(daoEnt) }

func (dao *daoEnt) InsertSpec(ctx context.Context, conn tx.Conn, spec termexp.DescSpec) error {
	rec, err := termexp.ConvertSpecToRec(spec)
	if err != nil {
		return err
	}
	dto, err := dataFromDescRec(rec)
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
			SetKind(int16(state.K)).
			SetSpec(jsontext.Value(raw)))
	}

	txCtx := entpgx.WithConn(ctx, conn.Pgx)
	err = dao.client.TermExp.CreateBulk(builders...).
		OnConflictColumns(termexpent.FieldID).
		DoNothing().
		Exec(txCtx)
	if err != nil {
		dao.log.Error("pool term Definition insertion failed", slog.Any("error", err))
		return err
	}

	dao.log.Log(ctx, lf.LevelTrace, "term Definition insertion succeeded", slog.Int("states", len(dto.States)))
	return nil
}

func (dao *daoEnt) SelectSpecByCK(ctx context.Context, conn tx.Conn, key contkey.ADT) (termexp.DescSpec, error) {
	rows, err := conn.Pgx.Query(ctx, selectTermExpTree, int64(key))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	states, err := pgx.CollectRows(rows, pgx.RowToStructByName[DescState])
	if err != nil {
		return nil, err
	}
	if len(states) == 0 {
		return nil, fmt.Errorf("pool term Definition %d does not exist", key)
	}
	rec, err := dataToDescRec(DescData{ExpCK: int64(key), States: states})
	if err != nil {
		return nil, err
	}
	return termexp.ConvertRecToSpec(rec), nil
}

func canonicalJSON(raw []byte) ([]byte, error) {
	value := jsontext.Value(raw)
	if err := value.Canonicalize(jsontext.CanonicalizeRawInts(false)); err != nil {
		return nil, err
	}
	return []byte(value), nil
}

const selectTermExpTree = `
WITH RECURSIVE exp_tree AS (
    SELECT top.*
    FROM pool_term_exps AS top
    WHERE top.exp_ck = $1
    UNION ALL
    SELECT sub.*
    FROM pool_term_exps AS sub
    JOIN exp_tree AS sup ON sub.sup_exp_ck = sup.exp_ck
)
SELECT exp_ck, sup_exp_ck, kind, spec
FROM exp_tree
`
