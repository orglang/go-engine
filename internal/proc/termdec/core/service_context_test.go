package core

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"orglang/go-engine/internal/adt/descroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/lib/tx"
)

func TestServiceInceptPropagatesCallerCancellationToRepository(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := &contextRepo{entered: make(chan struct{})}
	api := NewService(repo, contextDescRootRepo{}, nil,
		contextTransactor{}, slog.New(slog.NewTextHandler(io.Discard, nil)))

	result := make(chan error, 1)
	go func() {
		_, err := api.Incept(ctx, uniqsym.ADT{})
		result <- err
	}()

	<-repo.entered
	cancel()
	if err := <-result; !errors.Is(err, context.Canceled) {
		t.Fatalf("Incept error = %v, want context.Canceled", err)
	}
	if err := repo.ctxErr; !errors.Is(err, context.Canceled) {
		t.Fatalf("repository context error = %v, want context.Canceled", err)
	}
}

type contextRepo struct {
	entered chan struct{}
	ctxErr  error
}

func (r *contextRepo) AddRec(ctx context.Context, _ tx.Conn, _ DecRec) error {
	close(r.entered)
	<-ctx.Done()
	r.ctxErr = ctx.Err()
	return r.ctxErr
}
func (*contextRepo) GetRefs(context.Context, tx.Conn) ([]termsem.SemRef, error) { return nil, nil }
func (*contextRepo) GetQNByID(context.Context, tx.Conn, identity.ADT) (uniqsym.ADT, error) {
	return uniqsym.ADT{}, nil
}
func (*contextRepo) GetSnap(context.Context, tx.Conn, termsem.SemRef) (DecSnap, error) {
	return DecSnap{}, nil
}
func (*contextRepo) GetRecs(context.Context, tx.Conn, []identity.ADT) ([]DecRec, error) {
	return nil, nil
}
func (*contextRepo) SelectEnv(context.Context, tx.Conn, []identity.ADT) (map[identity.ADT]DecRec, error) {
	return nil, nil
}

type contextDescRootRepo struct{}

func (contextDescRootRepo) AddRef(context.Context, tx.Conn, descroot.SemRef) error   { return nil }
func (contextDescRootRepo) TouchRef(context.Context, tx.Conn, descroot.SemRef) error { return nil }

type contextTransactor struct{}

func (contextTransactor) ExplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	return fn(ctx, tx.Conn{})
}
func (contextTransactor) ImplicitTx(ctx context.Context, fn func(context.Context, tx.Conn) error) error {
	return fn(ctx, tx.Conn{})
}

var _ Repo = (*contextRepo)(nil)
var _ descroot.Repo = contextDescRootRepo{}
var _ tx.Transactor = contextTransactor{}
