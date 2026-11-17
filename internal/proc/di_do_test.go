package proc

import (
	"context"
	"slices"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/adt/typesem"
	webapi "orglang/go-engine/internal/app/web/api"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	compconn "orglang/go-engine/internal/proc/compconn/core"
	compexecadapter "orglang/go-engine/internal/proc/compexec/adapter"
	compexec "orglang/go-engine/internal/proc/compexec/core"
	termdecadapter "orglang/go-engine/internal/proc/termdec/adapter"
	termdec "orglang/go-engine/internal/proc/termdec/core"
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
		{"generated process execution controller", func() error {
			_, err := do.Invoke[compexecadapter.GeneratedController](injector)
			return err
		}},
		{"term-declaration controller", func() error {
			_, err := do.Invoke[webapi.ServerInterface](injector)
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

func TestRegisterRoutes(t *testing.T) {
	injector := newTestInjector()
	server := echo.New()

	if err := RegisterRoutes(server, injector); err != nil {
		t.Fatalf("register process routes: %v", err)
	}

	got := make([]string, 0, len(server.Routes()))
	for _, route := range server.Routes() {
		got = append(got, route.Method+" "+route.Path)
	}
	want := []string{
		"GET /api/v1/processes/executions/:id",
		"POST /api/v1/processes/executions/:id/steps",
		"POST /api/v1/processes/declarations",
		"GET /api/v1/processes/declarations/:id",
		"POST /api/v1/types",
		"GET /api/v1/types/:id",
		"PATCH /api/v1/types/:id",
		"POST /ssr/decs",
		"GET /ssr/decs",
		"GET /ssr/decs/:id",
		"POST /ssr/types",
		"GET /ssr/types",
		"GET /ssr/types/:id",
	}
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Fatalf("unexpected process routes: got %v, want %v", got, want)
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
