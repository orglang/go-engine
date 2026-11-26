package adapter

import (
	"context"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/uniqsym"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	"orglang/go-engine/internal/lib/tx"
)

type termDecRepo struct { dao *daoEnt }

func newTermDecRepo(dao *daoEnt) termdec.Repo { return termDecRepo{dao: dao} }

func (r termDecRepo) AddRec(ctx context.Context, conn tx.Conn, rec termdec.DecRec) error {
	return r.dao.AddRec(ctx, conn, termdec.DecRec{TermRef: rec.TermRef, TermQN: rec.TermQN, LiabVar: rec.LiabVar, AssetVars: rec.AssetVars})
}
func (r termDecRepo) GetRefs(ctx context.Context, conn tx.Conn) ([]termsem.SemRef, error) {
	return r.dao.GetRefs(ctx, conn)
}
func (r termDecRepo) GetSnap(ctx context.Context, conn tx.Conn, ref termsem.SemRef) (termdec.DecSnap, error) {
	return r.dao.GetSnap(ctx, conn, ref)
}
func (r termDecRepo) GetRecs(ctx context.Context, conn tx.Conn, ids []identity.ADT) ([]termdec.DecRec, error) {
	return r.dao.GetRecs(ctx, conn, ids)
}
func (r termDecRepo) SelectEnv(ctx context.Context, conn tx.Conn, ids []identity.ADT) (map[identity.ADT]termdec.DecRec, error) {
	return r.dao.SelectEnv(ctx, conn, ids)
}

func (r termDecRepo) GetQNByID(ctx context.Context, conn tx.Conn, id identity.ADT) (uniqsym.ADT, error) {
	return r.dao.GetQNByID(ctx, conn, id)
}
