package adapter

import (
	"context"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/termsem"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	termdef "orglang/go-engine/internal/proc/termdef/core"
	"orglang/go-engine/internal/lib/tx"
)

type termDecRepo struct { dao *daoEnt }

func newTermDecRepo(dao *daoEnt) termdec.Repo { return termDecRepo{dao: dao} }

func (r termDecRepo) AddRec(ctx context.Context, conn tx.Conn, rec termdec.DecRec) error {
	return r.dao.AddRec(ctx, conn, termdef.DecRec{TermRef: rec.TermRef, TermQN: rec.TermQN, LiabVar: rec.LiabVar, AssetVars: rec.AssetVars})
}
func (r termDecRepo) GetRefs(ctx context.Context, conn tx.Conn) ([]termsem.SemRef, error) {
	return r.dao.GetRefs(ctx, conn)
}
func (r termDecRepo) GetSnap(ctx context.Context, conn tx.Conn, ref termsem.SemRef) (termdec.DecSnap, error) {
	s, err := r.dao.GetSnap(ctx, conn, ref)
	return termdec.DecSnap{TermRef:s.TermRef, LiabVar:s.LiabVar, AssetVars:s.AssetVars}, err
}
func (r termDecRepo) GetRecs(ctx context.Context, conn tx.Conn, ids []identity.ADT) ([]termdec.DecRec, error) {
	recs, err := r.dao.GetRecs(ctx, conn, ids)
	if err != nil { return nil, err }
	out := make([]termdec.DecRec, len(recs))
	for i, rec := range recs { out[i] = termdec.DecRec{TermRef:rec.TermRef, TermQN:rec.TermQN, LiabVar:rec.LiabVar, AssetVars:rec.AssetVars} }
	return out, nil
}
func (r termDecRepo) SelectEnv(ctx context.Context, conn tx.Conn, ids []identity.ADT) (map[identity.ADT]termdec.DecRec, error) {
	recs, err := r.GetRecs(ctx, conn, ids)
	if err != nil { return nil, err }
	out := make(map[identity.ADT]termdec.DecRec, len(recs))
	for _, rec := range recs { out[rec.TermRef.TermID] = rec }
	return out, nil
}
