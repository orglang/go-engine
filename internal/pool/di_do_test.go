package pool

import (
	"context"
	"testing"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	compconn "orglang/go-engine/internal/pool/compconn/core"
	compexecadapter "orglang/go-engine/internal/pool/compexec/adapter"
	compexec "orglang/go-engine/internal/pool/compexec/core"
	compvar "orglang/go-engine/internal/pool/compvar/core"
	termdef "orglang/go-engine/internal/pool/termdef/core"
	typedef "orglang/go-engine/internal/pool/typedef/core"
	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

type testTransactor struct{}

func (testTransactor) ExplicitTx(context.Context, func(context.Context, tx.Conn) error) error {
	return nil
}

func (testTransactor) ImplicitTx(context.Context, func(context.Context, tx.Conn) error) error {
	return nil
}

func TestPackageProvidesPersistencePortsAndDefinitionServices(t *testing.T) {
	injector := do.New(
		lf.Package,
		do.Package(func(i do.Injector) {
			do.ProvideValue[tx.Transactor](i, testTransactor{})
		}),
		Package,
	)

	providers := []struct {
		name string
		call func() error
	}{
		{"pool connection repository", func() error {
			_, err := do.Invoke[compconn.Repo](injector)
			return err
		}},
		{"pool component-variable repository", func() error {
			_, err := do.Invoke[compvar.Repo](injector)
			return err
		}},
		{"pool computation-execution repository", func() error {
			_, err := do.Invoke[compexec.Repo](injector)
			return err
		}},
		{"pool type-expression repository", func() error {
			_, err := do.Invoke[typeexp.Repo](injector)
			return err
		}},
		{"pool type-definition repository", func() error {
			_, err := do.Invoke[typedef.Repo](injector)
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
		{"pool term-definition service", func() error {
			_, err := do.Invoke[termdef.API](injector)
			return err
		}},
		{"pool computation-execution service", func() error {
			_, err := do.Invoke[compexec.API](injector)
			return err
		}},
		{"pool computation-execution broker subscription", func() error {
			return compexecadapter.SubscribeBroker(injector)
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
