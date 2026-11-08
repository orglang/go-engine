# Architecture

This project follows a **screaming hexagonal architecture**: the repository structure should make its domain capabilities and use cases visible, while infrastructure details remain at the edges.

## Dependency direction

- Dependencies point inward. Domain and application code must not depend on database drivers, web frameworks, or concrete adapters.
- Core packages own the ports (interfaces) they need. Define repository and service contracts close to the consuming core package, commonly in `ds.go`.
- Adapters implement those ports and translate between infrastructure APIs and the project's contracts. Keep concrete driver dependencies in adapter packages.
- Inject dependencies through constructors; do not locate infrastructure from core code.
- Shared boundary contracts belong in `lib` when they are genuinely shared. Keep them small and independent of application-specific policy.

## Package layout

- `adt/` contains domain concepts and their persistence-facing ports.
- `proc/` and `pool/` group application/domain capabilities. Their `core` packages contain use-case logic; adjacent ports describe the capabilities core needs.
- `lib/` contains shared technical contracts and infrastructure integration. Concrete implementations belong in adapter packages.
- Package and directory names should describe the business capability first. Frameworks and storage technologies should appear only at the outer edges.

## Persistence and transactions

- Core code depends on `lib/tx` contracts, not directly on pgx or a pool implementation.
- Repositories receive the callback's `tx.Conn` and context explicitly. Do not acquire connections or begin transactions inside repositories.
- The transactor owns connection/transaction lifetime. Do not retain a `tx.Conn` after its callback returns.
- Use `ExplicitTx` when a use case needs atomic writes; use `ImplicitTx` for work that only needs a pooled connection.
- Nested transactor calls are not supported. Pass the existing callback context and connection through all repository calls in one unit of work.
- Keep pgx-specific behavior in the pgx adapter. Do not add driver-specific types to `proc/**/core` or `pool/**/core`.

## Change guidance

- Prefer business-focused changes in core and implement infrastructure changes behind existing ports.
- Add unit tests for core behavior using fake ports; test adapter behavior against the real infrastructure when lifecycle or driver semantics matter.
- Do not introduce interfaces or abstractions solely to make a test possible unless they represent a real boundary.
- Preserve the dependency direction when adding packages: outer layers may import inner layers, never the reverse.
