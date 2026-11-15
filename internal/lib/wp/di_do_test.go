package wp

import (
	"context"
	"testing"

	"github.com/alitto/pond/v2"
	"github.com/gammazero/workerpool"
	"github.com/samber/do/v2"

	"orglang/go-engine/internal/app"
)

func TestPackageProvidesPoolsAndStopsThemWithApp(t *testing.T) {
	application := app.NewApp(Package)

	pondPool, err := do.Invoke[pond.Pool](application.Container())
	if err != nil {
		t.Fatalf("resolve pond pool: %v", err)
	}
	workerPool, err := do.Invoke[*workerpool.WorkerPool](application.Container())
	if err != nil {
		t.Fatalf("resolve worker pool: %v", err)
	}

	if err := application.Stop(context.Background()); err != nil {
		t.Fatalf("stop application: %v", err)
	}
	if !pondPool.Stopped() {
		t.Fatal("expected pond pool to be stopped")
	}
	if !workerPool.Stopped() {
		t.Fatal("expected worker pool to be stopped")
	}
}
