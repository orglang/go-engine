package proc

import (
	"context"
	"testing"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/descsem"
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	compconn "orglang/go-engine/internal/proc/compconn/core"
	termdecadapter "orglang/go-engine/internal/proc/termdec/adapter"
	typedefadapter "orglang/go-engine/internal/proc/typedef/adapter"
	compexec "orglang/go-engine/internal/proc/compexec/core"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	termdef "orglang/go-engine/internal/proc/termdef/core"
	typedef "orglang/go-engine/internal/proc/typedef/core"
	typeexp "orglang/go-engine/internal/proc/typeexp/core"
)

type testTransactor struct{}

func (testTransactor) ExplicitTx(context.Context, func(context.Context, tx.Conn) error) error {
	return nil
}

func (testTransactor) ImplicitTx(context.Context, func(context.Context, tx.Conn) error) error {
	return nil
}

func TestPackageProvidesPersistencePorts(t *testing.T) {
	injector := newTestInjector()
	providers := []struct {
		name string
		call func() error
	}{
		{"connection repository", func() error {
			_, err := do.Invoke[compconn.Repo](injector)
			return err
		}},
		{"computation-execution repository", func() error {
			_, err := do.Invoke[compexec.Repo](injector)
			return err
		}},
		{"term-declaration repository", func() error {
			_, err := do.Invoke[termdec.Repo](injector)
			return err
		}},
		{"term-definition repository", func() error {
			_, err := do.Invoke[termdef.Repo](injector)
			return err
		}},
		{"type-definition repository", func() error {
			_, err := do.Invoke[typedef.Repo](injector)
			return err
		}},
		{"type-expression repository", func() error {
			_, err := do.Invoke[typeexp.Repo](injector)
			return err
		}},
		{"shared type-semantic repository", func() error {
			_, err := do.Invoke[typesem.Repo](injector)
			return err
		}},
		{"shared description-semantic repository", func() error {
			_, err := do.Invoke[descsem.Repo](injector)
			return err
		}},
	}

	for _, provider := range providers {
		t.Run(provider.name, func(t *testing.T) {
			if err := provider.call(); err != nil {
				t.Fatalf("resolve provider: %v", err)
			}
		})
	}
}

func TestPackageProvidesApplicationServices(t *testing.T) {
	injector := newTestInjector()

	providers := []struct {
		name string
		call func() error
	}{
		{"computation-execution service", func() error {
			_, err := do.Invoke[compexec.API](injector)
			return err
		}},
		{"term-declaration service", func() error {
			_, err := do.Invoke[termdec.API](injector)
			return err
		}},
		{"type-definition service", func() error {
			_, err := do.Invoke[typedef.API](injector)
			return err
		}},
	}

	for _, provider := range providers {
		t.Run(provider.name, func(t *testing.T) {
			if err := provider.call(); err != nil {
				t.Fatalf("resolve provider: %v", err)
			}
		})
	}
}

func TestPackageProvidesHTTPAdapters(t *testing.T) {
	injector := newTestInjector()

	providers := []struct {
		name string
		call func() error
	}{
		{"term-declaration controller", func() error {
			_, err := do.Invoke[termdecadapter.Controller](injector)
			return err
		}},
		{"term-declaration presenter", func() error {
			_, err := do.Invoke[termdecadapter.Presenter](injector)
			return err
		}},
		{"type-definition controller", func() error {
			_, err := do.Invoke[typedefadapter.Controller](injector)
			return err
		}},
		{"type-definition presenter", func() error {
			_, err := do.Invoke[typedefadapter.Presenter](injector)
			return err
		}},
	}

	for _, provider := range providers {
		t.Run(provider.name, func(t *testing.T) {
			if err := provider.call(); err != nil {
				t.Fatalf("resolve provider: %v", err)
			}
		})
	}
}

func newTestInjector() do.Injector {
	return do.New(
		lf.Package,
		do.Package(func(i do.Injector) {
			do.ProvideValue[tx.Transactor](i, testTransactor{})
		}),
		Package,
	)
}
