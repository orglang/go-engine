package app

import (
	"errors"
	"testing"

	"github.com/samber/do/v2"
)

type containerTestService struct {
	value string
}

func TestNewContainerResolvesPackageProviders(t *testing.T) {
	injector := NewContainer(do.Package(func(i do.Injector) {
		do.Provide[*containerTestService](i, func(do.Injector) (*containerTestService, error) {
			return &containerTestService{value: "provided"}, nil
		})
	}))

	got, err := do.Invoke[*containerTestService](injector)
	if err != nil {
		t.Fatalf("invoke service: %v", err)
	}
	if got.value != "provided" {
		t.Fatalf("service value = %q, want %q", got.value, "provided")
	}
}

func TestNewContainerReturnsResolutionErrors(t *testing.T) {
	injector := NewContainer()

	if _, err := do.Invoke[*containerTestService](injector); err == nil {
		t.Fatal("expected missing provider error")
	}

	providerErr := errors.New("provider failed")
	injector = NewContainer(do.Package(func(i do.Injector) {
		do.Provide[*containerTestService](i, func(do.Injector) (*containerTestService, error) {
			return nil, providerErr
		})
	}))

	if _, err := do.Invoke[*containerTestService](injector); !errors.Is(err, providerErr) {
		t.Fatalf("invoke error = %v, want wrapped %v", err, providerErr)
	}
}

func TestNewContainerOverridesAreIsolated(t *testing.T) {
	packageWithProvider := do.Package(func(i do.Injector) {
		do.Provide[*containerTestService](i, func(do.Injector) (*containerTestService, error) {
			return &containerTestService{value: "default"}, nil
		})
	})

	first := NewContainer(packageWithProvider)
	second := NewContainer(packageWithProvider)
	do.OverrideValue[*containerTestService](first, &containerTestService{value: "test override"})

	firstService, err := do.Invoke[*containerTestService](first)
	if err != nil {
		t.Fatalf("invoke overridden service: %v", err)
	}
	secondService, err := do.Invoke[*containerTestService](second)
	if err != nil {
		t.Fatalf("invoke independent service: %v", err)
	}

	if firstService.value != "test override" {
		t.Fatalf("first container value = %q, want %q", firstService.value, "test override")
	}
	if secondService.value != "default" {
		t.Fatalf("second container value = %q, want %q", secondService.value, "default")
	}
}
