package pool

import (
	"testing"

	"github.com/samber/do/v2"

	"orglang/go-engine/internal/lib/lf"
	compconn "orglang/go-engine/internal/pool/compconn/core"
	compvar "orglang/go-engine/internal/pool/compvar/core"
	typeexp "orglang/go-engine/internal/pool/typeexp/core"
)

func TestPackageProvidesPersistencePorts(t *testing.T) {
	injector := do.New(lf.Package, Package)

	if _, err := do.Invoke[compconn.Repo](injector); err != nil {
		t.Fatalf("invoke pool connection repository: %v", err)
	}
	if _, err := do.Invoke[compvar.Repo](injector); err != nil {
		t.Fatalf("invoke pool component-variable repository: %v", err)
	}
	if _, err := do.Invoke[typeexp.Repo](injector); err != nil {
		t.Fatalf("invoke pool type-expression repository: %v", err)
	}
}
