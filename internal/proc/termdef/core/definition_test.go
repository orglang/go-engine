package core

import (
	"context"
	"errors"
	"testing"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

type definitionTestTransactor struct{}

func (definitionTestTransactor) ExplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	return fn(ctx, tx.Conn{})
}
func (definitionTestTransactor) ImplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	return fn(ctx, tx.Conn{})
}

type definitionTestRepo struct {
	rec  DefRec
	spec ExpSpec
	err  error
	key  contkey.ADT
}

func (r *definitionTestRepo) SetExpCK(context.Context, tx.Conn, identity.ADT, contkey.ADT) error {
	return nil
}
func (r *definitionTestRepo) GetRecByID(context.Context, tx.Conn, identity.ADT) (DefRec, error) {
	return r.rec, r.err
}
func (r *definitionTestRepo) InsertSpec(context.Context, tx.Conn, ExpSpec) error { return nil }
func (r *definitionTestRepo) SelectSpecByCK(_ context.Context, _ tx.Conn, key contkey.ADT) (ExpSpec, error) {
	r.key = key
	return r.spec, nil
}

type definitionTestDecRepo struct{ qn uniqsym.ADT }

func (*definitionTestDecRepo) AddRec(context.Context, tx.Conn, termdec.DecRec) error { return nil }
func (*definitionTestDecRepo) GetIDs(context.Context, tx.Conn) ([]identity.ADT, error) {
	return nil, nil
}
func (r *definitionTestDecRepo) GetQNByID(context.Context, tx.Conn, identity.ADT) (uniqsym.ADT, error) {
	return r.qn, nil
}
func (*definitionTestDecRepo) GetSnap(context.Context, tx.Conn, identity.ADT) (termdec.DecSnap, error) {
	return termdec.DecSnap{}, nil
}
func (*definitionTestDecRepo) GetRecs(context.Context, tx.Conn, []identity.ADT) ([]termdec.DecRec, error) {
	return nil, nil
}
func (*definitionTestDecRepo) SelectEnv(context.Context, tx.Conn, []identity.ADT) (map[identity.ADT]termdec.DecRec, error) {
	return nil, nil
}

func TestRetrieveDefinition(t *testing.T) {
	termID := identity.New()
	wantKey := contkey.ADT(42)
	wantQN, err := uniqsym.ConvertFromString("example.term")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("returns current definition", func(t *testing.T) {
		repo := &definitionTestRepo{rec: DefRec{TermID: termID, ExpCK: wantKey}, spec: ExpCloseSpec{}}
		decRepo := &definitionTestDecRepo{qn: wantQN}
		service := NewService(decRepo, repo, definitionTestTransactor{})
		got, err := service.RetrieveDefinition(context.Background(), termID)
		if err != nil {
			t.Fatalf("RetrieveDefinition() error = %v", err)
		}
		if got.TermQN != wantQN {
			t.Fatalf("RetrieveDefinition().TermQN = %v, want %v", got.TermQN, wantQN)
		}
		if _, ok := got.Exp.(ExpCloseSpec); !ok {
			t.Fatalf("RetrieveDefinition().Exp = %T, want ExpCloseSpec", got.Exp)
		}
		if repo.key != wantKey {
			t.Fatalf("SelectSpecByCK() key = %v, want %v", repo.key, wantKey)
		}
	})
	t.Run("returns missing term error", func(t *testing.T) {
		repo := &definitionTestRepo{err: ErrTermMissing(termID)}
		service := NewService(&definitionTestDecRepo{}, repo, definitionTestTransactor{})
		_, err := service.RetrieveDefinition(context.Background(), termID)
		if !errors.Is(err, ErrTermDoesNotExist) {
			t.Fatalf("RetrieveDefinition() error = %v, want ErrTermDoesNotExist", err)
		}
	})
	t.Run("returns missing definition error", func(t *testing.T) {
		repo := &definitionTestRepo{rec: DefRec{TermID: termID}}
		service := NewService(&definitionTestDecRepo{}, repo, definitionTestTransactor{})
		_, err := service.RetrieveDefinition(context.Background(), termID)
		if !errors.Is(err, ErrDefinitionMissing) {
			t.Fatalf("RetrieveDefinition() error = %v, want ErrDefinitionMissing", err)
		}
	})
}
