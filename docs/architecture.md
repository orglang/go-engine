# Architecture

## Dependency direction

The target architecture is ports-and-adapters, with dependencies pointing inward:

- **Domain** contains business concepts and invariants. It must not depend on transport, persistence, frameworks, or concrete adapters.
- **Application** coordinates use cases. It depends on domain types and ports, never on infrastructure implementations.
- **Ports** are interfaces owned by the domain/application side and describe the capabilities those layers need.
- **Adapters** implement ports or translate external input into application calls. Database drivers, HTTP handlers, SDK integrations, and framework-specific code belong here.
- **Composition roots** are executable entry points that construct concrete adapters and wire them to application services.

Prefer cohesive packages and narrow public APIs. Keep interfaces at the consumer boundary; introduce abstractions only where they separate a real dependency.

## Direction of migration

The repository currently has a mixed layout (app, pool, proc, adt, lib, and db). Some pool and proc packages already distinguish core from adapter, but the separation is incomplete: core packages still reach into infrastructure packages such as lib/db.

Migrate incrementally rather than moving directories mechanically:

1. Identify each infrastructure dependency used by domain/application code.
2. Define the smallest port at the consuming boundary.
3. Move the implementation behind an adapter and wire it in the composition root.
4. Move packages only when their responsibility and dependency direction are clear.
5. Preserve behavior and public interfaces unless a compatibility change is explicitly documented.

A possible end-state is:

```text
cmd/                    # executable entry points / composition roots
internal/
  domain/               # domain model and invariants
  application/          # use cases
  ports/                # inbound and outbound contracts
  adapters/
    inbound/            # HTTP, CLI, message consumers
    outbound/           # PostgreSQL, SDKs, external services
```

This is a guide, not a mandatory template: retain cohesive package boundaries where the existing domain model calls for them.

## Enforced boundaries

The depguard configuration applies to Go files under packages named core or domain. It currently prevents those packages from importing concrete adapter packages and the wire-format SDK.

The repository still has existing direct infrastructure dependencies from core packages (notably lib/db). Those are migration work, not exemptions to the target architecture. Extend the depguard deny rules as each dependency is moved behind a port; do not add new infrastructure imports to core/domain packages.
