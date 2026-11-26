package proc

import (
	"context"
	"testing"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/typesem"
	webapi "orglang/go-engine/internal/app/web/api"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	compconn "orglang/go-engine/internal/proc/compconn/core"
	compexecadapter "orglang/go-engine/internal/proc/compexec/adapter"
	compexec "orglang/go-engine/internal/proc/compexec/core"
	termdec "orglang/go-engine/internal/proc/termdec/core"
	termdefadapter "orglang/go-engine/internal/proc/termdef/adapter"
	termdef "orglang/go-engine/internal/proc/termdef/core"
	typedefadapter "orglang/go-engine/internal/proc/typedef/adapter"
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
		{"process description-semantic repository", func() error {
			_, err := do.Invoke[typedef.DescSemRepo](injector)
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
		{"unified term lifecycle service", func() error {
			_, err := do.Invoke[termdef.API](injector)
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
		{"generated process execution controller", func() error {
			_, err := do.Invoke[compexecadapter.GeneratedController](injector)
			return err
		}},
		{"term-definition controller", func() error {
			_, err := do.Invoke[webapi.ServerInterface](injector)
			return err
		}},
		{"term-definition presenter", func() error {
			_, err := do.Invoke[termdefadapter.Presenter](injector)
			return err
		}},
		{"generated type-definition controller", func() error {
			_, err := do.Invoke[typedefadapter.GeneratedController](injector)
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
