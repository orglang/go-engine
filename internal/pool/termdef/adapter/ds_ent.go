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

func NewDaoEnt(client *ent.Client, log *slog.Logger) termdef.Repo {
	name := slog.String("name", reflect.TypeFor[daoEnt]().Name())
	return &daoEnt{client: client, log: log.With(name)}
}

var _ termdef.Repo = (*daoEnt)(nil)

func (dao *daoEnt) GetByTermID(ctx context.Context, conn tx.Conn, termID identity.ADT) (termdef.DefRec, error) {
	if termID.IsEmpty() {
		return termdef.DefRec{}, identity.ErrEmpty
	}
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
	if expCK == 0 {
		return termdef.DefRec{TermID: termID}, nil
	}
	exp, err := dao.selectSpecByCK(ctx, conn, expCK)
	if err != nil {
		return termdef.DefRec{}, err
	}
	return termdef.DefRec{TermID: termID, Exp: exp}, nil
}
