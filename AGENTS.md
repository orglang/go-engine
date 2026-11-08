# Architecture

This project follows a **screaming hexagonal architecture**: package structure exposes business capabilities, while infrastructure stays at the edges.

- Dependencies point inward. Domain and application code must not depend on drivers, frameworks, or concrete adapters.
- Core packages own the ports they need. Adapters implement those ports; inject implementations through constructors.
- Keep shared technical contracts in `lib/`. Keep concrete infrastructure integrations in adapter packages.
- Organize `adt/`, `proc/`, and `pool/` around domain concepts and business capabilities. Keep technology-specific names at the outer edges.
- Core code depends on `lib/tx`, never directly on pgx. Repositories receive `tx.Conn` and context from the transaction callback; they do not acquire connections or start transactions.
- The transactor owns connection and transaction lifetime. Use `ExplicitTx` for atomic writes and `ImplicitTx` when only a pooled connection is needed. Nested transactions are unsupported.
