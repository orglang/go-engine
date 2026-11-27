package pool

import (
	"context"
	"testing"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	execconn "orglang/go-engine/internal/pool/execconn/core"
	compvar "orglang/go-engine/internal/pool/compvar/core"
	execevaladapter "orglang/go-engine/internal/pool/execeval/adapter"
	execeval "orglang/go-engine/internal/pool/execeval/core"
	termdecadapter "orglang/go-engine/internal/pool/termdec/adapter"
	termdec "orglang/go-engine/internal/pool/termdec/core"
	termdefadapter "orglang/go-engine/internal/pool/termdef/adapter"
	termdef "orglang/go-engine/internal/pool/termdef/core"
	typedefadapter "orglang/go-engine/internal/pool/typedef/adapter"
	typedef "orglang/go-engine/internal/pool/typedef/core"
	"orglang/go-engine/internal/proc"
)

type testTransactor struct{}

func (testTransactor) ExplicitTx(context.Context, func(context.Context, tx.Conn) error) error {
	return nil
}

func (testTransactor) ImplicitTx(context.Context, func(context.Context, tx.Conn) error) error {
	return nil
}

func TestPackageProvidesPersistencePortsAndDefinitionServices(t *testing.T) {
	injector := newTestInjector()

	providers := []struct {
		name string
		call func() error
	}{
		{"pool connection repository", func() error {
			_, err := do.Invoke[execconn.Repo](injector)
			return err
		}},
		{"pool component-variable repository", func() error {
			_, err := do.Invoke[compvar.Repo](injector)
			return err
		}},
		{"pool computation-execution repository", func() error {
			_, err := do.Invoke[execeval.Repo](injector)
			return err
		}},
		{"pool type-definition repository", func() error {
			_, err := do.Invoke[typedef.Repo](injector)
			return err
		}},
		{"pool term-declaration repository", func() error {
			_, err := do.Invoke[termdec.Repo](injector)
			return err
		}},
		{"pool term-definition repository", func() error {
			_, err := do.Invoke[termdef.Repo](injector)
			return err
		}},
		{"pool type-definition service", func() error {
			_, err := do.Invoke[typedef.API](injector)
			return err
		}},
		{"pool term-declaration service", func() error {
			_, err := do.Invoke[termdec.API](injector)
			return err
		}},
		{"pool term-definition service", func() error {
			_, err := do.Invoke[termdef.API](injector)
			return err
		}},
		{"pool computation-execution service", func() error {
			_, err := do.Invoke[execeval.API](injector)
			return err
		}},
		{"pool computation-execution broker subscription", func() error {
			return execevaladapter.SubscribeBroker(injector)
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

func TestPackageProvidesHTTPControllers(t *testing.T) {
	injector := newTestInjector()

	providers := []struct {
		name string
		call func() error
	}{
		{"computation-execution controller", func() error {
			_, err := do.Invoke[execevaladapter.Controller](injector)
			return err
		}},
		{"generated type-definition controller", func() error {
			_, err := do.Invoke[typedefadapter.GeneratedController](injector)
			return err
		}},
		{"generated term-declaration controller", func() error {
			_, err := do.Invoke[termdecadapter.GeneratedController](injector)
			return err
		}},
		{"generated term-definition controller", func() error {
			_, err := do.Invoke[termdefadapter.GeneratedController](injector)
			return err
		}},
	}

	for _, provider := range providers {
		t.Run(provider.name, func(t *testing.T) {
			if err := provider.call(); err != nil {
				t.Fatalf("resolve controller: %v", err)
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
		proc.Package,
		Package,
	)
}
