package poolcomm

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/poolcfg"
	"orglang/go-engine/adt/poolctx"
	"orglang/go-engine/adt/poolenv"
)

type Repo interface {
	InsertRec(db.Source, CommRec) error
	InsertRecs(db.Source, []CommRec) error
	SelectEnvSnapByEnvSpec(db.Source, poolenv.EnvSpec) (poolenv.EnvSnap, error)
	SelectCtxSnapByCtxSpec(db.Source, poolctx.CtxSpec) (poolctx.CtxSnap, error)
	SelectCfgSnapBySpec(db.Source, poolcfg.CfgSpec) (poolcfg.CfgSnap, error)
	SelectRecByRef(db.Source, commsem.SemRef) (CommRec, error)
}

type commRecDS struct {
	ChnlID string
	ChnlON int64
	K      commKind
}

type commKind uint8

const (
	unkComm commKind = iota
	Publication
	Subscription
)
