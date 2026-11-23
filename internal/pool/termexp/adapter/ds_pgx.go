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

func (dao *daoPgx) InsertSpec(ctx context.Context, conn tx.Conn, spec termexp.ExpSpec) (contkey.ADT, error) {
	nodes, err := termexp.ContentNodes(spec)
	if err != nil { return contkey.Zero, err }
	for _, node := range nodes {
		_, err := conn.Pgx.Exec(ctx, insertTermExp,
			int64(node.ExpCK), int64(node.SupExpCK), node.Kind, string(node.Spec))
		if err != nil {
			dao.log.Error("term expression insertion failed", slog.Int64("exp_ck", int64(node.ExpCK)), slog.Any("error", err))
			return contkey.Zero, err
		}
		var existing expNodeRow
		err = conn.Pgx.QueryRow(ctx, selectTermExpByCK, int64(node.ExpCK)).Scan(
			&existing.SupExpCK, &existing.Kind, &existing.Spec)
		if err != nil { return contkey.Zero, err }
		canonical, err := canonicalStoredJSON(existing.Spec)
		if err != nil { return contkey.Zero, err }
		if existing.SupExpCK != int64(node.SupExpCK) || existing.Kind != node.Kind || !bytes.Equal(canonical, node.Spec) {
			return contkey.Zero, fmt.Errorf("pool term expression content-key collision or parent conflict for %d", node.ExpCK)
		}
	}
	return nodes[0].ExpCK, nil
}

func (dao *daoPgx) SelectSpecByCK(ctx context.Context, conn tx.Conn, key contkey.ADT) (termexp.ExpSpec, error) {
	rows, err := conn.Pgx.Query(ctx, selectTermExpTree, int64(key))
	if err != nil { return nil, err }
	defer rows.Close()
	dtos, err := pgx.CollectRows(rows, pgx.RowToStructByName[expNodeRow])
	if err != nil { return nil, err }
	if len(dtos) == 0 { return nil, fmt.Errorf("pool term expression %d does not exist", key) }
	nodes := make([]termexp.ExpNode, len(dtos))
	for i, dto := range dtos {
		expCK, convErr := contkey.ConvertFromInt(dto.ExpCK); if convErr != nil { return nil, convErr }
		supCK, convErr := contkey.ConvertFromInt(dto.SupExpCK); if convErr != nil { return nil, convErr }
		nodes[i] = termexp.ExpNode{ExpCK: expCK, SupExpCK: supCK, Kind: dto.Kind, Spec: dto.Spec}
	}
	return termexp.BuildSpecFromNodes(nodes, key)
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
