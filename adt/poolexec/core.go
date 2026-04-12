package poolexec

import (
	"context"
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/adt/poolcomm"
	"orglang/go-engine/adt/pooldec"
	"orglang/go-engine/adt/poolstep"
	"orglang/go-engine/adt/poolvar"
	"orglang/go-engine/adt/procdec"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Run(ExecSpec) (implsem.SemRef, error) // aka Create
}

type ExecSpec struct {
	// ссылка на декларацию вновь создаваемого пула
	DescQN uniqsym.ADT
	// внутренняя и внешняя ссылки на вновь создаваемый пул
	LiabVar implvar.VarSpec
	// внутренние и внешние ссылки на ранее созданные пулы
	AssetVars []implvar.VarSpec
}

type ExecRec struct {
	ImplRef  implsem.SemRef
	LiabMode implvar.Mode
}

type ExecSnap struct {
	ImplRef implsem.SemRef
	LiabVar implvar.VarRec
}

type service struct {
	poolExecRepo Repo
	implSemRepo  implsem.Repo
	commSemRepo  commsem.Repo
	poolDecRepo  pooldec.Repo
	poolVarRepo  poolvar.Repo
	poolCommRepo poolcomm.Repo
	operator     db.Operator
	log          *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func newService(
	poolExecRepo Repo,
	implSemRepo implsem.Repo,
	commSemRepo commsem.Repo,
	poolDecRepo pooldec.Repo,
	poolVarRepo poolvar.Repo,
	poolCommRepo poolcomm.Repo,
	poolStepRepo poolstep.Repo,
	xactDefRepo xactexp.Repo,
	xactExpRepo xactexp.Repo,
	procDecRepo procdec.Repo,
	operator db.Operator,
	log *slog.Logger,
) *service {
	name := slog.String("name", reflect.TypeFor[service]().Name())
	return &service{
		poolExecRepo, implSemRepo, commSemRepo, poolDecRepo, poolVarRepo, poolCommRepo,
		operator, log.With(name),
	}
}

func (s *service) Run(spec ExecSpec) (_ implsem.SemRef, err error) {
	ctx := context.Background()
	specAttr := slog.Any("spec", spec)
	s.log.Debug("creation started", specAttr)
	var poolDec pooldec.DecRec
	getErr1 := s.operator.Implicit(ctx, func(ds db.Source) error {
		poolDec, err = s.poolDecRepo.GetRecByQN(ds, spec.DescQN)
		return err
	})
	if getErr1 != nil {
		s.log.Error("creation failed", specAttr)
		return implsem.SemRef{}, getErr1
	}
	assetQNs := make([]uniqsym.ADT, 0, len(spec.AssetVars))
	for _, assetVar := range spec.AssetVars {
		if assetVar.ImplQN == spec.LiabVar.ImplQN {
			continue
		}
		assetQNs = append(assetQNs, assetVar.ImplQN)
	}
	var assetExecs map[uniqsym.ADT]ExecSnap
	getErr2 := s.operator.Implicit(ctx, func(ds db.Source) error {
		assetExecs, err = s.poolExecRepo.GetSnapMapByQNs(ds, assetQNs)
		return err
	})
	if getErr2 != nil {
		s.log.Error("creation failed", specAttr)
		return implsem.SemRef{}, getErr2
	}
	newImplSem := implsem.SemRec{ImplRef: implsem.NewRef(), ImplQN: spec.LiabVar.ImplQN, Kind: implsem.PoolKind}
	newCommSem := commsem.SemRec{CommRef: commsem.NewRef(), Kind: commsem.Pool}
	newConn := poolcomm.ConnRec{CommRef: newCommSem.CommRef, CommON: 0}
	newExec := ExecRec{ImplRef: newImplSem.ImplRef, LiabMode: implvar.StructMode}
	newLiabVar := implvar.StructRec{
		ImplRef: newImplSem.ImplRef,
		CommRef: newCommSem.CommRef,
		ChnlID:  identity.New(),
		ChnlPH:  spec.LiabVar.ChnlPH,
		ChnlBS:  implvar.LiabSide,
		ExpVK:   poolDec.LiabVar.ExpVK,
	}
	newAssetVars := make([]implvar.VarRec, 0, len(spec.AssetVars)+1)
	for _, assetVar := range spec.AssetVars {
		var commRef commsem.SemRef
		var chnlID identity.ADT
		var expVK valkey.ADT
		assetExec, ok := assetExecs[assetVar.ImplQN]
		if !ok && assetVar.ImplQN == spec.LiabVar.ImplQN {
			commRef = newCommSem.CommRef
			chnlID = newLiabVar.ChnlID
			expVK = newLiabVar.ExpVK
		} else {
			commRef = assetExec.LiabVar.GetCommRef()
			chnlID = assetExec.LiabVar.GetChnlID()
			expVK = assetExec.LiabVar.GetExpVK()
		}
		newAssetVars = append(newAssetVars, implvar.StructRec{
			ImplRef: newImplSem.ImplRef,
			CommRef: commRef,
			ChnlID:  chnlID,
			ChnlPH:  assetVar.ChnlPH,
			ChnlBS:  implvar.AssetSide,
			ExpVK:   expVK,
		})
	}
	transactErr := s.operator.Explicit(ctx, func(ds db.Source) error {
		err = s.implSemRepo.AddRec(ds, newImplSem)
		if err != nil {
			return err
		}
		err = s.poolExecRepo.AddRec(ds, newExec)
		if err != nil {
			return err
		}
		err = s.poolVarRepo.AddRecs(ds, append(newAssetVars, newLiabVar))
		if err != nil {
			return err
		}
		err = s.commSemRepo.AddRec(ds, newCommSem)
		if err != nil {
			return err
		}
		return s.poolCommRepo.AddRec(ds, newConn)
	})
	if transactErr != nil {
		s.log.Error("creation failed", specAttr)
		return implsem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", slog.Any("ref", newImplSem.ImplRef))
	return newImplSem.ImplRef, nil
}
