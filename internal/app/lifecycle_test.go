package app

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/samber/do/v2"
)

func TestLifecycleStartsInOrderAndStopsInReverse(t *testing.T) {
	var calls []string
	lifecycle := NewLifecycle()

	for _, name := range []string{"db", "server"} {
		name := name
		if err := lifecycle.Append(Hook{
			OnStart: func(context.Context) error {
				calls = append(calls, "start:"+name)
				return nil
			},
			OnStop: func(context.Context) error {
				calls = append(calls, "stop:"+name)
				return nil
			},
		}); err != nil {
			t.Fatalf("append %s hook: %v", name, err)
		}
	}

	if err := lifecycle.Start(context.Background()); err != nil {
		t.Fatalf("start lifecycle: %v", err)
	}
	if err := lifecycle.Stop(context.Background()); err != nil {
		t.Fatalf("stop lifecycle: %v", err)
	}

	want := []string{"start:db", "start:server", "stop:server", "stop:db"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("lifecycle calls = %v, want %v", calls, want)
	}
}

func TestLifecycleRollsBackStartedHooksWhenStartupFails(t *testing.T) {
	startupErr := errors.New("server failed")
	var calls []string
	lifecycle := NewLifecycle()

	_ = lifecycle.Append(Hook{
		OnStart: func(context.Context) error {
			calls = append(calls, "start:db")
			return nil
		},
		OnStop: func(context.Context) error {
			calls = append(calls, "stop:db")
			return nil
		},
	})
	_ = lifecycle.Append(Hook{
		OnStart: func(context.Context) error {
			calls = append(calls, "start:server")
			return startupErr
		},
		OnStop: func(context.Context) error {
			calls = append(calls, "stop:server")
			return nil
		},
	})

	if err := lifecycle.Start(context.Background()); !errors.Is(err, startupErr) {
		t.Fatalf("start error = %v, want %v", err, startupErr)
	}

	want := []string{"start:db", "start:server", "stop:db"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("lifecycle calls = %v, want %v", calls, want)
	}
	if err := lifecycle.Stop(context.Background()); err != nil {
		t.Fatalf("stop after failed startup: %v", err)
	}
}

type lifecycleShutdownService struct {
	called bool
}

func (s *lifecycleShutdownService) Shutdown(context.Context) error {
	s.called = true
	return nil
}

func TestAppRegistersLifecycleBeforeApplicationPackages(t *testing.T) {
	var packageLifecycle Lifecycle
	var resolveErr error
	service := &lifecycleShutdownService{}
	app := NewApp(do.Package(func(i do.Injector) {
		packageLifecycle, resolveErr = do.Invoke[Lifecycle](i)
		do.ProvideValue[*lifecycleShutdownService](i, service)
	}))

	if resolveErr != nil {
		t.Fatalf("resolve lifecycle during package registration: %v", resolveErr)
	}
	if packageLifecycle == nil {
		t.Fatal("package resolved a nil lifecycle")
	}

	gotLifecycle, err := do.Invoke[Lifecycle](app.Container())
	if err != nil {
		t.Fatalf("resolve lifecycle: %v", err)
	}
	if gotLifecycle != packageLifecycle {
		t.Fatal("package and application resolved different lifecycle instances")
	}

	if _, err := do.Invoke[*lifecycleShutdownService](app.Container()); err != nil {
		t.Fatalf("resolve shutdown service: %v", err)
	}
	if err := app.Start(context.Background()); err != nil {
		t.Fatalf("start app: %v", err)
	}
	if err := app.Stop(context.Background()); err != nil {
		t.Fatalf("stop app: %v", err)
	}
	if !service.called {
		t.Fatal("expected container shutdown to call service Shutdown")
	}
}

func TestLifecycleRejectsLateHooksAndRepeatedStart(t *testing.T) {
	lifecycle := NewLifecycle()
	if err := lifecycle.Append(Hook{OnStart: func(context.Context) error { return nil }}); err != nil {
		t.Fatalf("append hook: %v", err)
	}
	if err := lifecycle.Start(context.Background()); err != nil {
		t.Fatalf("start lifecycle: %v", err)
	}
	if err := lifecycle.Append(Hook{OnStop: func(context.Context) error { return nil }}); err == nil {
		t.Fatal("expected append after startup to fail")
	}
	if err := lifecycle.Start(context.Background()); err == nil {
		t.Fatal("expected repeated start to fail")
	}
	if err := lifecycle.Stop(context.Background()); err != nil {
		t.Fatalf("stop lifecycle: %v", err)
	}
	if err := lifecycle.Stop(context.Background()); err != nil {
		t.Fatalf("repeated stop: %v", err)
	}
}
