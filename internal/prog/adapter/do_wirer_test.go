package adapter

import (
	"slices"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/lib/lf"
	coreprog "orglang/go-engine/internal/prog/core"
)

func TestPackageProvidesProgramServiceAndController(t *testing.T) {
	injector := do.New(lf.Package, Package)

	if _, err := do.Invoke[coreprog.API](injector); err != nil {
		t.Fatalf("resolve program service: %v", err)
	}
	controller, err := do.Invoke[Controller](injector)
	if err != nil {
		t.Fatalf("resolve program controller: %v", err)
	}

	server := echo.New()
	if err := RegisterRoutes(server, controller); err != nil {
		t.Fatalf("register program routes: %v", err)
	}

	got := make([]string, 0, len(server.Router().Routes()))
	for _, route := range server.Router().Routes() {
		got = append(got, route.Method+" "+route.Path)
	}
	want := []string{"POST " + programRoute}
	slices.Sort(got)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Fatalf("registered program routes: got %v, want %v", got, want)
	}
}
