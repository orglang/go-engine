package poolstep

import (
	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/poolcfg"
	"orglang/go-engine/adt/poolctx"
	"orglang/go-engine/adt/poolenv"
)

type Repo interface {
	InsertRec(db.Source, StepRec) error
	AddRecs(db.Source, []StepRec) error
	SelectEnvSnapByEnvSpec(db.Source, poolenv.EnvSpec) (poolenv.EnvSnap, error)
	SelectCtxSnapByCtxSpec(db.Source, poolctx.CtxQry) (poolctx.CtxSnap, error)
	SelectCfgSnapBySpec(db.Source, poolcfg.CfgSpec) (poolcfg.CfgSnap, error)
	SelectRecByRef(db.Source, implsem.SemRef) (StepRec, error)
}

type commRecDS struct {
	ChnlID string
	ChnlON int64
	K      commKind
}

type commKind uint16

const (
	unkComm commKind = iota
	Publication
	Subscription
)
