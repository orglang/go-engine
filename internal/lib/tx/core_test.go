package tx

import (
	"context"
	"errors"
	"testing"
)

type testTransactor struct {
	explicit func(context.Context, func(context.Context, Conn) error) error
	implicit func(context.Context, func(context.Context, Conn) error) error
}

func (t testTransactor) ExplicitTx(ctx context.Context, fn func(context.Context, Conn) error) error {
	return t.explicit(ctx, fn)
}

func (t testTransactor) ImplicitTx(ctx context.Context, fn func(context.Context, Conn) error) error {
	return t.implicit(ctx, fn)
}

func TestExplicitReturnsCallbackValue(t *testing.T) {
	want := "result"
	transactor := testTransactor{
		explicit: func(ctx context.Context, fn func(context.Context, Conn) error) error {
			return fn(ctx, Conn{})
		},
	}

	got, err := Explicit(context.Background(), transactor, func(context.Context, Conn) (string, error) {
		return want, nil
	})
	if err != nil {
		t.Fatalf("Explicit() error = %v", err)
	}
	if got != want {
		t.Fatalf("Explicit() = %q, want %q", got, want)
	}
}

func TestImplicitReturnsZeroValueOnCallbackError(t *testing.T) {
	wantErr := errors.New("callback failed")
	transactor := testTransactor{
		implicit: func(ctx context.Context, fn func(context.Context, Conn) error) error {
			return fn(ctx, Conn{})
		},
	}

	got, err := Implicit(context.Background(), transactor, func(context.Context, Conn) (string, error) {
		return "partial result", wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Implicit() error = %v, want %v", err, wantErr)
	}
	if got != "" {
		t.Fatalf("Implicit() = %q, want zero value", got)
	}
}

func TestGenericHelpersRejectNilCallbacks(t *testing.T) {
	transactor := testTransactor{}

	if _, err := Explicit[string](context.Background(), transactor, nil); err == nil {
		t.Fatal("Explicit() error = nil, want non-nil")
	}
	if _, err := Implicit[string](context.Background(), transactor, nil); err == nil {
		t.Fatal("Implicit() error = nil, want non-nil")
	}
}
