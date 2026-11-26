package adapter

import (
	"context"

	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
	termdec "orglang/go-engine/internal/pool/termdec/core"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

type termDecRepo struct { dao *daoEnt }

func newTermDecRepo(dao *daoEnt) termdec.Repo { return termDecRepo{dao: dao} }

func (r termDecRepo) AddRec(ctx context.Context, conn tx.Conn, rec termdec.DecRec) error {
	return r.dao.AddRec(ctx, conn, termdef.DecRec{TermRef:rec.TermRef, LiabVar:rec.LiabVar, AssetVars:rec.AssetVars})
}
func (r termDecRepo) GetRecByQN(ctx context.Context, conn tx.Conn, qn uniqsym.ADT) (termdec.DecRec, error) {
	rec, err := r.dao.GetRecByQN(ctx, conn, qn)
	return termdec.DecRec{TermRef:rec.TermRef, LiabVar:rec.LiabVar, AssetVars:rec.AssetVars}, err
}
