package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/proc/ent"
	termdef "orglang/go-engine/internal/proc/termdef/core"
)

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
