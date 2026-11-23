package adapter

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/lib/tx"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

type daoPgx struct { log *slog.Logger }

type expNodeRow struct {
	ExpCK    int64  `db:"exp_ck"`
	SupExpCK int64  `db:"sup_exp_ck"`
	Kind     int16  `db:"kind"`
	Spec     []byte `db:"spec"`
}

func newRepo() termexp.Repo { return new(daoPgx) }

func NewDaoPgx(log *slog.Logger) termexp.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{log: log.With(name)}
}

func (dao *daoPgx) InsertSpec(ctx context.Context, conn tx.Conn, spec termexp.DescSpec) (contkey.ADT, error) {
	rec, err := termexp.ConvertSpecToRec(spec)
	if err != nil { return contkey.Zero, err }
	dto, err := dataFromDescRec(rec)
	if err != nil { return contkey.Zero, err }
	for _, state := range dto.States {
		_, err := conn.Pgx.Exec(ctx, insertTermExp,
			state.ExpCK, state.SupExpCK, int16(state.K), string(state.Spec))
		if err != nil {
			dao.log.Error("pool term Definition insertion failed", slog.Int64("exp_ck", state.ExpCK), slog.Any("error", err))
			return contkey.Zero, err
		}
		var existing expNodeRow
		err = conn.Pgx.QueryRow(ctx, selectTermExpByCK, state.ExpCK).Scan(
			&existing.SupExpCK, &existing.Kind, &existing.Spec)
		if err != nil { return contkey.Zero, err }
		canonical, err := canonicalStoredJSON(existing.Spec)
		if err != nil { return contkey.Zero, err }
		expected, err := canonicalStoredJSON(state.Spec)
		if err != nil { return contkey.Zero, err }
		if existing.SupExpCK != state.SupExpCK || existing.Kind != int16(state.K) || !bytes.Equal(canonical, expected) {
			return contkey.Zero, fmt.Errorf("pool term Definition content-key collision or parent conflict for %d", state.ExpCK)
		}
	}
	return contkey.ConvertFromInt(dto.ExpCK)
}

func (dao *daoPgx) SelectSpecByCK(ctx context.Context, conn tx.Conn, key contkey.ADT) (termexp.DescSpec, error) {
	rows, err := conn.Pgx.Query(ctx, selectTermExpTree, int64(key))
	if err != nil { return nil, err }
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[expNodeRow])
	if err != nil { return nil, err }
	if len(dtos) == 0 { return nil, fmt.Errorf("pool term Definition %d does not exist", key) }
	data := DescData{ExpCK: int64(key), States: make([]DescState, len(dtos))}
	for i, row := range dtos {
		data.States[i] = DescState{ExpCK:row.ExpCK, SupExpCK:row.SupExpCK, K:ExpKind(row.Kind), Spec:row.Spec}
	}
	rec, err := dataToDescRec(data)
	if err != nil { return nil, err }
	return termexp.ConvertRecToSpec(rec), nil
}

func canonicalStoredJSON(raw []byte) ([]byte, error) {
	value := jsontext.Value(raw)
	if err := value.Canonicalize(jsontext.CanonicalizeRawInts(false)); err != nil { return nil, err }
	return []byte(value), nil
}

const insertTermExp = `
INSERT INTO pool_term_exps (exp_ck, sup_exp_ck, kind, spec)
VALUES ($1, $2, $3, $4::jsonb)
ON CONFLICT (exp_ck) DO NOTHING
`

const selectTermExpByCK = `
SELECT sup_exp_ck, kind, spec
FROM pool_term_exps
WHERE exp_ck = $1
`

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
