package bootstrap

import (
	"context"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/app"
	"orglang/go-engine/internal/app/web"
	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"
	"orglang/go-engine/internal/pool"
	"orglang/go-engine/internal/proc"
	progadapter "orglang/go-engine/internal/prog/adapter"
)

type testTransactor struct{}

func (testTransactor) ExplicitTx(context.Context, func(context.Context, tx.Conn) error) error {
	return nil
}

func (testTransactor) ImplicitTx(context.Context, func(context.Context, tx.Conn) error) error {
	return nil
}

func TestNewWithPackagesComposesHTTPRoutes(t *testing.T) {
	infrastructure := do.Package(func(i do.Injector) {
		do.ProvideValue[tx.Transactor](i, testTransactor{})
		do.ProvideValue[*echo.Echo](i, echo.New())
	})

	application, err := NewWithPackages(
		lf.Package,
		infrastructure,
		pool.Package,
		proc.Package,
		progadapter.Package,
		web.Package,
	)
	if err != nil {
		t.Fatalf("compose application: %v", err)
	}
	t.Cleanup(func() {
		if err := application.Stop(context.Background()); err != nil {
			t.Errorf("stop application: %v", err)
		}
	})

	server, err := do.Invoke[*echo.Echo](application.Container())
	if err != nil {
		t.Fatalf("resolve Echo server: %v", err)
	}

	routes := make(map[string]struct{}, len(server.Routes()))
	for _, route := range server.Routes() {
		routes[route.Method+" "+route.Path] = struct{}{}
	}

	required := []string{
		"GET /",
		"GET /api/v1/processes/executions/:id",
		"POST /api/v1/processes/executions/:id/steps",
		"POST /api/v1/pools/execs",
		"POST /api/v1/processes/declarations",
		"POST /api/v1/types",
		"POST /api/v1/programs",
	}
	for _, route := range required {
		if _, ok := routes[route]; !ok {
			t.Errorf("application route %q was not registered", route)
		}
	}

	legacy := []string{
		"GET /api/v1/procs/:id",
		"POST /api/v1/procs/:id/steps",
	}
	for _, route := range legacy {
		if _, ok := routes[route]; ok {
			t.Errorf("legacy application route %q is still registered", route)
		}
	}

	if len(routes) != len(server.Routes()) {
		t.Fatal("application registered duplicate HTTP routes")
	}

	if _, err := do.Invoke[app.Lifecycle](application.Container()); err != nil {
		t.Fatalf("resolve application lifecycle: %v", err)
	}
}
