package lf

import (
	"log/slog"
	"testing"

	"github.com/samber/do/v2"
)

func TestPackageProvidesLogger(t *testing.T) {
	injector := do.New(Package)

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		t.Fatalf("invoke logger: %v", err)
	}
	if logger == nil {
		t.Fatal("logger is nil")
	}
}
