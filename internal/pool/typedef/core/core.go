package core

import (
	"context"
	"fmt"
	"log/slog"

	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
)

type API interface {
	Create(context.Context, DefSpec) (identity.ADT, descroot.SemRef, error)
}

type DefSpec struct {
	TypeQN  uniqsym.ADT
	TypeExp ExpSpec
}

type DefRec struct {
	TypeQN  uniqsym.ADT
	TypeID  identity.ADT
	DescRef descroot.SemRef
	ExpCK   contkey.ADT
	TypeExp ExpRec
}

type DefSnap struct {
	TypeID  identity.ADT
	DescRef descroot.SemRef
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

func (s *service) Create(ctx context.Context, spec DefSpec) (_ identity.ADT, _ descroot.SemRef, err error) {
	qnAttr := slog.Any("qn", spec.TypeQN)
	s.log.Debug("creation started", qnAttr, slog.Any("spec", spec))
	newExp, convErr := ConvertSpecToRec(spec.TypeExp)
	if convErr != nil {
		return identity.ADT{}, descroot.SemRef{}, convErr
	}
	typeID := identity.New()
	descRef := descroot.New()
	newDef := DefRec{TypeID: typeID, DescRef: descRef, TypeQN: spec.TypeQN, ExpCK: newExp.Key(), TypeExp: newExp}
	transactErr := s.transactor.ExplicitTx(ctx, func(txCtx context.Context, conn tx.Conn) error {
		if err := s.descRootRepo.AddRef(txCtx, conn, descRef); err != nil {
			return err
		}
		return s.typeDefRepo.AddRec(txCtx, conn, newDef)
	})
	if transactErr != nil {
		s.log.Error("creation failed", qnAttr)
		return identity.ADT{}, descroot.SemRef{}, transactErr
	}
	s.log.Debug("creation succeed", qnAttr, slog.Any("typeID", typeID), slog.Any("descRef", descRef))
	return typeID, descRef, nil
}

func ErrMissingInEnv(want contkey.ADT) error {
	return fmt.Errorf("exp missing in env: %v", want)
}
