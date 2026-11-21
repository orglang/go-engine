# Architecture

This project follows a **screaming hexagonal architecture**: package structure exposes business capabilities, while infrastructure stays at the edges.

- Dependencies point inward. Domain and application code must not depend on drivers, frameworks, or concrete adapters.
- Core packages own the ports they need. Adapters implement those ports; inject implementations through constructors.
- Keep shared technical contracts in `lib/`. Keep concrete infrastructure integrations in adapter packages.
- Organize `adt/`, `proc/`, and `pool/` around domain concepts and business capabilities. Keep technology-specific names at the outer edges.
- Core code depends on `lib/tx`, never directly on pgx. Repositories receive `tx.Conn` and context from the transaction callback; they do not acquire connections or start transactions.
- The transactor owns connection and transaction lifetime. Use `ExplicitTx` for atomic writes and `ImplicitTx` when only a pooled connection is needed. Nested transactions are unsupported.

## SQL formatting

- Write SQL in Go raw string literals with the first SQL token at column 0, independent of the surrounding Go indentation.
- Put the opening and closing backticks on their own lines for multiline SQL strings.
- Format SQL keywords in uppercase and indent nested clauses and expressions with four spaces.
- Keep SQL structure readable; do not indent SQL to match the Go source indentation.
- Do not add runtime trimming or rewriting of query strings to compensate for source formatting.
