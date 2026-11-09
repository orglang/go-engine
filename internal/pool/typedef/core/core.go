package core

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"

	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

type API interface {
	Create(DefSpec) (typesem.SemRef, error)
}

type DefSpec struct {
	TypeQN  uniqsym.ADT
	TypeExp typeexp.ExpSpec
}

type DefRec struct {
	TypeRef typesem.SemRef
	ExpCK   contkey.ADT
}

type DefSnap struct {
	TypeRef typesem.SemRef
	TypeExp typeexp.ExpSpec
}

type service struct {
	typeDefRepo Repo
	typeExpRepo typeexp.Repo
	descSemRepo descsem.Repo
	transactor  tx.Transactor
	log         *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func NewService(
	typeDefRepo Repo,
	typeExpRepo typeexp.Repo,
	descSemRepo descsem.Repo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	return &service{typeDefRepo, typeExpRepo, descSemRepo, transactor, log}
}

func (s *service) Create(spec DefSpec) (_ typesem.SemRef, err error) {
	ctx := context.Background()
	qnAttr := slog.Any("qn", spec.TypeQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newExp, convErr := typeexp.ConvertSpecToRec(spec.TypeExp)
	if convErr != nil {
		return typesem.SemRef{}, convErr
	}
	newDef := DefRec{TypeRef: typesem.New(), ExpCK: newExp.Key()}
	newBind := descsem.SemRec{DescQN: spec.TypeQN, DescID: newDef.TypeRef.TypeID, Kind: descsem.TypeKind}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		err = s.descSemRepo.AddRec(txCtx, conn, newBind)
		if err != nil {
			return err
		}
		err = s.typeExpRepo.AddRec(txCtx, conn, newExp, newDef.TypeRef)
		if err != nil {
			return err
		}
		return s.typeDefRepo.AddRec(txCtx, conn, newDef)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return typesem.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("ref", newDef.TypeRef))
	return newDef.TypeRef, nil
}

func errOptimisticUpdate(got seqnum.ADT) error {
	return fmt.Errorf("entity concurrent modification: got revision %v", got)
}

func ErrMissingInEnv(want contkey.ADT) error {
	return fmt.Errorf("exp missing in env: %v", want)
}
