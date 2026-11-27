package core

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
)

type API interface {
	Create(context.Context, DefSpec) (typesem.SemRef, error)
}

type DefSpec struct {
	TypeQN  uniqsym.ADT
	TypeExp ExpSpec
}

type DefRec struct {
	TypeQN  uniqsym.ADT
	TypeRef typesem.SemRef
	ExpCK   contkey.ADT
	TypeExp ExpRec
}

type DefSnap struct {
	TypeRef typesem.SemRef
	TypeExp ExpSpec
}

type service struct {
	typeDefRepo  Repo
	descRootRepo descroot.Repo
	transactor   tx.Transactor
	log          *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func NewService(
	typeDefRepo Repo,
	descRootRepo descroot.Repo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	return &service{typeDefRepo: typeDefRepo, descRootRepo: descRootRepo, transactor: transactor, log: log}
}

func (s *service) Create(ctx context.Context, spec DefSpec) (_ typesem.SemRef, err error) {
	qnAttr := slog.Any("qn", spec.TypeQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newExp, convErr := ConvertSpecToRec(spec.TypeExp)
	if convErr != nil {
		return typesem.SemRef{}, convErr
	}
	newDef := DefRec{TypeQN: spec.TypeQN, TypeRef: typesem.New(), ExpCK: newExp.Key(), TypeExp: newExp}
	newRoot := descroot.SemRef{DescID: newDef.TypeRef.DescID, DescRN: newDef.TypeRef.DescRN}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		err = s.descRootRepo.AddRef(txCtx, conn, newRoot)
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

func ErrMissingInEnv(want contkey.ADT) error {
	return fmt.Errorf("exp missing in env: %v", want)
}
