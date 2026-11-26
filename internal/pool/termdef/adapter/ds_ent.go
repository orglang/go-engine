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
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool/ent"
	termdef "orglang/go-engine/internal/pool/termdef/core"
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
func newRepo() termdef.Repo { return new(daoEnt) }

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
SELECT def.term_id, COALESCE(def.exp_ck, 0) AS exp_ck
FROM pool_desc_terms AS def
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
	expCK, err := contkey.ConvertFromInt(dto.ExpCK)
	if err != nil {
		return termdef.DefRec{}, err
	}
	return termdef.DefRec{TermID: termID, ExpCK: expCK}, nil
}
