# Architecture

## Organize around the domain

The repository follows a screaming architecture: its top-level packages and
directories should make the product's major capabilities recognizable. The
existing `proc` and `pool` modules are meaningful domain boundaries and
should remain recognizable rather than being flattened into repository-wide
technical layers.

Use Go's `internal/` mechanism where it improves encapsulation. A possible
direction for the main modules is:

```text
internal/
  proc/
    core/                 # process domain and application behavior
    adapter/              # infrastructure and external-system integrations
    ...                   # module-specific ports and supporting packages
  pool/
    core/                 # pool domain and application behavior
    adapter/              # infrastructure and external-system integrations
    ...                   # module-specific ports and supporting packages
  ...                     # other cohesive internal capabilities
cmd/                      # executable entry points / composition roots
```

This is illustrative, not a mandatory directory template. Preserve the existing
internal organization of `proc/` and `pool/` where it is cohesive. Move code
under `internal/` when it should not be imported by external consumers; do not
move packages solely to make the tree match this example. Keep other top-level
packages when they represent intentional public APIs or independently useful
capabilities.

## Hexagonal boundaries inside modules

Hexagonal architecture describes dependency direction, not a required
repository-wide directory hierarchy:

- **Core** owns domain rules and application behavior for its module. It must
  not depend on concrete adapters, transport, persistence drivers, frameworks,
  or wire-format SDKs.
- **Ports** are small interfaces defined at the consuming boundary. Keep them
  close to the module and behavior that needs them; avoid a central repository-
  wide ports package unless a contract is genuinely shared.
- **Adapters** implement outbound ports or translate inbound external input
  into calls to the module. Database access, SDK integrations, transport
  handlers, and framework-specific code belong outside the core.
- **Composition roots** construct concrete adapters and wire them to the
  relevant module. Keep wiring at executable entry points or other explicit
  application boundaries, not inside domain logic.

Dependencies point inward: adapters may depend on core-owned contracts, while
core must not import adapters. Prefer narrow interfaces and introduce
abstractions only where they separate a real dependency.

## Direction of migration

The repository currently has a mixed layout (app, pool, proc, adt, lib, and db).
Some `pool` and `proc` packages already distinguish core from adapter, but
the separation is incomplete: core packages still reach into infrastructure
packages such as `lib/db`.

Migrate incrementally rather than moving directories mechanically:

1. Identify each infrastructure dependency used by a module's core.
2. Define the smallest port at the consuming module boundary.
3. Move or adapt the implementation behind an adapter.
4. Wire the adapter from the appropriate composition root.
5. Add or tighten architecture checks once the dependency has actually been
   removed.
6. Preserve behavior and public interfaces unless a compatibility change is
   explicitly documented.

Keep `proc` and `pool` recognizable throughout the migration. Do not introduce
global `domain/`, `application/`, `ports/`, or `adapters/` trees merely to
classify packages by technical layer.

## Enforced boundaries

The depguard configuration applies to Go files under packages named `core` or
`domain`. It currently prevents those packages from importing concrete
adapter packages and the wire-format SDK.

The repository still has existing direct infrastructure dependencies from
core packages (notably `lib/db`). Those are migration work, not exemptions to
the target architecture. Extend the depguard deny rules as each dependency is
moved behind a port; do not add new infrastructure imports to core/domain
packages.
