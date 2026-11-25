package core

import (
	"context"
	"errors"
	"testing"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
)

type definitionTestTransactor struct{}

func (definitionTestTransactor) ExplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	return fn(ctx, tx.Conn{})
}

func (definitionTestTransactor) ImplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	return fn(ctx, tx.Conn{})
}

type definitionTestDefRepo struct {
	rec DefRec
	err error
	qn  uniqsym.ADT
}

func (r definitionTestDefRepo) SetExpCK(context.Context, tx.Conn, identity.ADT, contkey.ADT) error {
	return nil
}

func (r definitionTestDefRepo) GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error) {
	return r.rec, r.err
}

func (r definitionTestDefRepo) GetQNByID(context.Context, tx.Conn, identity.ADT) (uniqsym.ADT, error) {
	return r.qn, nil
}

type definitionTestTermExpRepo struct {
	spec ExpSpec
	key  contkey.ADT
}

func (r *definitionTestTermExpRepo) InsertSpec(context.Context, tx.Conn, ExpSpec) error {
	return nil
}

func (r *definitionTestTermExpRepo) SelectSpecByCK(_ context.Context, _ tx.Conn, key contkey.ADT) (ExpSpec, error) {
	r.key = key
	return r.spec, nil
}

func TestRetrieveDefinition(t *testing.T) {
	termID := identity.New()
	wantKey := contkey.ADT(42)
	wantQN := uniqsym.ADT("example.term")

	t.Run("returns current definition", func(t *testing.T) {
		termExpRepo := &definitionTestTermExpRepo{spec: ExpReleaseSpec{}}
		service := NewService(nil, definitionTestDefRepo{rec: DefRec{TermID: termID, ExpCK: wantKey}, qn: wantQN}, termExpRepo, nil, nil, nil, definitionTestTransactor{}, nil)

		got, err := service.RetrieveDefinition(context.Background(), termID)
		if err != nil {
			t.Fatalf("RetrieveDefinition() error = %v", err)
		}
		if got.TermQN != wantQN {
			t.Fatalf("RetrieveDefinition().TermQN = %v, want %v", got.TermQN, wantQN)
		}
		if _, ok := got.Exp.(ExpReleaseSpec); !ok {
			t.Fatalf("RetrieveDefinition().Exp = %T, want ExpReleaseSpec", got.Exp)
		}
		if termExpRepo.key != wantKey {
			t.Fatalf("SelectSpecByCK() key = %v, want %v", termExpRepo.key, wantKey)
		}
	})

	t.Run("returns missing term error", func(t *testing.T) {
		service := NewService(nil, definitionTestDefRepo{err: ErrTermMissing(termID)}, &definitionTestTermExpRepo{}, nil, nil, nil, definitionTestTransactor{}, nil)

		_, err := service.RetrieveDefinition(context.Background(), termID)
		if !errors.Is(err, ErrTermDoesNotExist) {
			t.Fatalf("RetrieveDefinition() error = %v, want ErrTermDoesNotExist", err)
		}
	})

	t.Run("returns missing definition error", func(t *testing.T) {
		service := NewService(nil, definitionTestDefRepo{rec: DefRec{TermID: termID}}, &definitionTestTermExpRepo{}, nil, nil, nil, definitionTestTransactor{}, nil)

		_, err := service.RetrieveDefinition(context.Background(), termID)
		if !errors.Is(err, ErrDefinitionMissing) {
			t.Fatalf("RetrieveDefinition() error = %v, want ErrDefinitionMissing", err)
		}
	})
}
