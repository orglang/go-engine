package poolexec

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/semcomm"
	"orglang/go-engine/adt/semterm"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
	"orglang/go-engine/pool/commexch"
	"orglang/go-engine/pool/commturn"
	compvar1 "orglang/go-engine/pool/compvar"
	termdec1 "orglang/go-engine/pool/termdec"
	"orglang/go-engine/pool/typeexp"
	"orglang/go-engine/proc/termdec"
)

type API interface {
	Run(ExecSpec) (semterm.TermRef, error) // aka Create
}

type ExecSpec struct {
	// ссылка на декларацию вновь создаваемого пула
	TermQN uniqsym.ADT
	// внутренняя и внешняя ссылки на вновь создаваемый пул
	LiabVar compvar.VarSpec
	// внутренние и внешние ссылки на ранее созданные пулы
	AssetVars []compvar.VarSpec
}

type ExecRec struct {
	CompRef  semterm.TermRef
	TermQN   uniqsym.ADT
	LiabMode compvar.Mode
}

type ExecSnap struct {
	CompRef semterm.TermRef
	LiabVar compvar.VarRec
}

type service struct {
	poolExecRepo Repo
	commSemRepo  semcomm.Repo
	poolDecRepo  termdec1.Repo
	poolVarRepo  compvar1.Repo
	poolCommRepo commexch.Repo
	operator     db.Operator
	log          *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	poolExecRepo Repo,
	commSemRepo semcomm.Repo,
	poolDecRepo termdec1.Repo,
	poolVarRepo compvar1.Repo,
	poolCommRepo commexch.Repo,
	poolStepRepo commturn.Repo,
	xactDefRepo typeexp.Repo,
	xactExpRepo typeexp.Repo,
	procDecRepo termdec.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{
		poolExecRepo, commSemRepo, poolDecRepo, poolVarRepo, poolCommRepo,
		operator, log.With(name),
	}
}

func (s *service) Run(spec ExecSpec) (_ semterm.TermRef, err error) {
	ctx := context.Background()
	specAttr := slog.Any("spec", spec)
	s.log.Debug("creation started", specAttr)
	var poolDec termdec1.DecRec
	getErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
		poolDec, err = s.poolDecRepo.GetRecByQN(ds, spec.TermQN)
		return err
	})
	if getErr1 != nil {
		s.log.Error("creation failed", specAttr)
		return semterm.TermRef{}, getErr1
	}
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars))
	for _, assetVar := range spec.AssetVars {
		if assetVar.TermQN == spec.LiabVar.TermQN {
			continue
		}
		assetQNs = append(assetQNs, assetVar.TermQN)
	}
	var assetExecs map[uniqsym.ADT]ExecSnap
	getErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
		assetExecs, err = s.poolExecRepo.GetSnapMapByQNs(ds, assetQNs)
		return err
	})
	if getErr2 != nil {
		s.log.Error("creation failed", specAttr)
		return semterm.TermRef{}, getErr2
	}
	newExec := ExecRec{CompRef: semterm.NewRef(), TermQN: spec.LiabVar.TermQN, LiabMode: compvar.StructMode}
	newComm := commexch.ExchRec{CommRef: semcomm.NewRef(), OffsetNr: 0}
	newLiabVar := compvar.StructRec{
		CompRef: newExec.CompRef,
		CommRef: newComm.CommRef,
		ChnlID:  identity.New(),
		ChnlPH:  spec.LiabVar.ChnlPH,
		ChnlBS:  compvar.LiabSide,
		ExpVK:   poolDec.LiabVar.ExpVK,
	}
	newAssetVars := make([]compvar.VarRec, 0, len(spec.AssetVars)+1)
	for _, assetVar := range spec.AssetVars {
		var commRef semcomm.CommRef
		var chnlID identity.ADT
		var expVK valkey.ADT
		assetExec, ok := assetExecs[assetVar.TermQN]
		if !ok && assetVar.TermQN == spec.LiabVar.TermQN {
			commRef = newComm.CommRef
			chnlID = newLiabVar.ChnlID
			expVK = newLiabVar.ExpVK
		} else {
			commRef = assetExec.LiabVar.GetCommRef()
			chnlID = assetExec.LiabVar.GetChnlID()
			expVK = assetExec.LiabVar.GetExpVK()
		}
		newAssetVars = append(newAssetVars, compvar.StructRec{
			CompRef: newExec.CompRef,
			CommRef: commRef,
			ChnlID:  chnlID,
			ChnlPH:  assetVar.ChnlPH,
			ChnlBS:  compvar.AssetSide,
			ExpVK:   expVK,
		})
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.poolExecRepo.AddRec(ds, newExec)
		if err != nil {
			return err
		}
		err = s.poolVarRepo.AddRecs(ds, append(newAssetVars, newLiabVar))
		if err != nil {
			return err
		}
		return s.poolCommRepo.AddRec(ds, newComm)
	})
	if transactErr != nil {
		s.log.Error("creation failed", specAttr)
		return semterm.TermRef{}, transactErr
	}
	s.log.Debug("creation succeed", slog.Any("ref", newExec.CompRef))
	return newExec.CompRef, nil
}
