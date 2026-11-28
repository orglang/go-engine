package kv

import (
	"os"
	"testing"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/lib/lf"
)

func TestPackageProvidesConfigurationLoader(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile("reference.yaml", []byte("service:\n  name: test\n"), 0o600); err != nil {
		t.Fatalf("write reference config: %v", err)
	}

	injector := do.New(lf.Package, Package)
	loader, err := do.Invoke[Loader](injector)
	if err != nil {
		t.Fatalf("invoke loader: %v", err)
	}

	var got struct {
		Name string `mapstructure:"name"`
	}
	if err := loader.Load("service", &got); err != nil {
		t.Fatalf("load service config: %v", err)
	}
	if got.Name != "test" {
		t.Fatalf("service name = %q, want %q", got.Name, "test")
	}
}

func TestPackageReturnsMissingReferenceConfigError(t *testing.T) {
	t.Chdir(t.TempDir())

	injector := do.New(lf.Package, Package)
	if _, err := do.Invoke[Loader](injector); err == nil {
		t.Fatal("expected missing reference config error")
	}
}
