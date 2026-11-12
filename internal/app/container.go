package app

import "github.com/samber/do/v2"

// NewContainer creates an isolated dependency-injection container and applies
// the supplied packages. Callers own the returned container and its shutdown.
func NewContainer(packages ...func(do.Injector)) do.Injector {
	return do.New(packages...)
}
