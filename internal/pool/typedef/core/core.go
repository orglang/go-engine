package core

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/progroot"
	"orglang/go-engine/internal/adt/uniqsym"
)

type API interface {
	Create(context.Context, DefSpec) (identity.ADT, progroot.RootRef, error)
}

type DefSpec struct {
	TypeQN  uniqsym.ADT
	TypeExp ExpSpec
}

type DefRec struct {
	TypeQN  uniqsym.ADT
	TypeID  identity.ADT
	ProgRef progroot.RootRef
	ExpCK   contkey.ADT
	TypeExp ExpRec
}

type DefSnap struct {
	TypeID  identity.ADT
	ProgRef progroot.RootRef
	TypeExp ExpSpec
}

type service struct {
	typeDefRepo  Repo
	progRootRepo progroot.Repo
	transactor   tx.Transactor
	log          *slog.Logger
}

// for compilation purposes
func newAPI() API {
	return new(service)
}

func NewService(
	typeDefRepo Repo,
	progRootRepo progroot.Repo,
	transactor tx.Transactor,
	log *slog.Logger,
) API {
	return &service{typeDefRepo: typeDefRepo, progRootRepo: progRootRepo, transactor: transactor, log: log}
}

func (s *service) Create(ctx context.Context, spec DefSpec) (_ identity.ADT, _ progroot.RootRef, err error) {
	qnAttr := slog.Any("qn", spec.TypeQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newExp, convErr := ConvertSpecToRec(spec.TypeExp)
	if convErr != nil {
		return identity.ADT{}, progroot.RootRef{}, convErr
	}
	typeID := identity.New()
	progRef := progroot.New()
	newDef := DefRec{TypeID: typeID, ProgRef: progRef, TypeQN: spec.TypeQN, ExpCK: newExp.Key(), TypeExp: newExp}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if err := s.progRootRepo.AddRef(txCtx, conn, progRef); err != nil {
			return err
		}
		return s.typeDefRepo.AddRec(txCtx, conn, newDef)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return identity.ADT{}, progroot.RootRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("typeID", typeID), slog.Any("progRef", progRef))
	return typeID, progRef, nil
}

func ErrMissingInEnv(want contkey.ADT) error {
	return fmt.Errorf("exp missing in env: %v", want)
}
