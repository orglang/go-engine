package core

import (
	"context"
	"errors"
	"testing"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/lib/tx"
	termexp "orglang/go-engine/internal/proc/termexp/core"
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
}

func (r definitionTestDefRepo) SetExpCK(context.Context, tx.Conn, identity.ADT, contkey.ADT) error {
	return nil
}

func (r definitionTestDefRepo) GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error) {
	return r.rec, r.err
}

type definitionTestTermExpRepo struct {
	spec termexp.DescSpec
	key  contkey.ADT
}

func (r *definitionTestTermExpRepo) Insert(context.Context, tx.Conn, termexp.CompRec) error {
	return nil
}

func (r *definitionTestTermExpRepo) InsertSpec(context.Context, tx.Conn, termexp.DescSpec) error {
	return nil
}

func (r *definitionTestTermExpRepo) SelectSpecByCK(_ context.Context, _ tx.Conn, key contkey.ADT) (termexp.DescSpec, error) {
	r.key = key
	return r.spec, nil
}

func TestRetrieveDefinition(t *testing.T) {
	termID := identity.New()
	wantKey := contkey.ADT(42)

	t.Run("returns current definition", func(t *testing.T) {
		termExpRepo := &definitionTestTermExpRepo{spec: termexp.DescCloseSpec{}}
		service := NewService(nil, definitionTestDefRepo{rec: DefRec{ExpCK: wantKey}}, termExpRepo, nil, nil, nil, definitionTestTransactor{}, nil)

		got, err := service.RetrieveDefinition(context.Background(), termID)
		if err != nil {
			t.Fatalf("RetrieveDefinition() error = %v", err)
		}
		if _, ok := got.(termexp.DescCloseSpec); !ok {
			t.Fatalf("RetrieveDefinition() = %T, want termexp.DescCloseSpec", got)
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
		service := NewService(nil, definitionTestDefRepo{rec: DefRec{}}, &definitionTestTermExpRepo{}, nil, nil, nil, definitionTestTransactor{}, nil)

		_, err := service.RetrieveDefinition(context.Background(), termID)
		if !errors.Is(err, ErrDefinitionMissing) {
			t.Fatalf("RetrieveDefinition() error = %v, want ErrDefinitionMissing", err)
		}
	})
}
