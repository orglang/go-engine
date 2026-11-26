package adapter

import (
	"context"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
	termdec "orglang/go-engine/internal/pool/termdec/core"
	)

type termDecRepo struct { dao *daoEnt }

func newTermDecRepo(dao *daoEnt) termdec.Repo { return termDecRepo{dao: dao} }

func (r termDecRepo) AddRec(ctx context.Context, conn tx.Conn, rec termdec.DecRec) error {
	return r.dao.AddRec(ctx, conn, termdec.DecRec{TermRef:rec.TermRef, LiabVar:rec.LiabVar, AssetVars:rec.AssetVars})
}
func (r termDecRepo) GetRecByQN(ctx context.Context, conn tx.Conn, qn uniqsym.ADT) (termdec.DecRec, error) {
	return r.dao.GetRecByQN(ctx, conn, qn)
}

func (r termDecRepo) GetQNByID(ctx context.Context, conn tx.Conn, id identity.ADT) (uniqsym.ADT, error) {
	return r.dao.GetQNByID(ctx, conn, id)
}
