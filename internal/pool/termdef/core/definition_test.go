package core

import (
	"context"
	"errors"
	"testing"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
	termdec "orglang/go-engine/internal/pool/termdec/core"
)

type definitionTestTransactor struct{}

func (definitionTestTransactor) ExplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	return fn(ctx, tx.Conn{})
}

func (definitionTestTransactor) ImplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	return fn(ctx, tx.Conn{})
}

type definitionTestRepo struct {
	rec DefRec
	err error
}

func (r *definitionTestRepo) GetByTermID(context.Context, tx.Conn, identity.ADT) (DefRec, error) {
	return r.rec, r.err
}

type definitionTestDecRepo struct {
	qn uniqsym.ADT
}

func (definitionTestDecRepo) AddRec(context.Context, tx.Conn, termdec.DecRec) error {
	return nil
}

func (definitionTestDecRepo) GetRecByQN(context.Context, tx.Conn, uniqsym.ADT) (termdec.DecRec, error) {
	return termdec.DecRec{}, nil
}

func (r definitionTestDecRepo) GetQNByID(context.Context, tx.Conn, identity.ADT) (uniqsym.ADT, error) {
	return r.qn, nil
}

func TestRetrieveDefinition(t *testing.T) {
	termID := identity.New()
	wantQN, err := uniqsym.ConvertFromString("example.term")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("returns current definition", func(t *testing.T) {
		repo := &definitionTestRepo{rec: DefRec{TermID: termID, Exp: ExpReleaseSpec{}}}
		service := NewService(definitionTestDecRepo{qn: wantQN}, repo, definitionTestTransactor{})

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
	})

	t.Run("returns missing term error", func(t *testing.T) {
		repo := &definitionTestRepo{err: ErrTermMissing(termID)}
		service := NewService(definitionTestDecRepo{}, repo, definitionTestTransactor{})

		_, err := service.RetrieveDefinition(context.Background(), termID)
		if !errors.Is(err, ErrTermDoesNotExist) {
			t.Fatalf("RetrieveDefinition() error = %v, want ErrTermDoesNotExist", err)
		}
	})

	t.Run("returns missing definition error", func(t *testing.T) {
		repo := &definitionTestRepo{rec: DefRec{TermID: termID}}
		service := NewService(definitionTestDecRepo{}, repo, definitionTestTransactor{})

		_, err := service.RetrieveDefinition(context.Background(), termID)
		if !errors.Is(err, ErrDefinitionMissing) {
			t.Fatalf("RetrieveDefinition() error = %v, want ErrDefinitionMissing", err)
		}
	})
}
