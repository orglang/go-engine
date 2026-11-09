# ADR-0003: Content Key terminology and persistence compatibility

- **Status:** Accepted
- **Date:** 2026-11-09

## Context

The domain concept previously named Value Key is a content-derived identity. The old terminology is misleading, but the key's calculation, composition, inversion, and lifecycle are not changing.

The Go package and domain symbols are implementation details. PostgreSQL columns, JSON properties, and SDK fields are persistence or externally consumed representations, so each boundary needs an explicit compatibility decision.

## Decision

Use **Content Key** as the canonical domain terminology in Go:

- package `internal/adt/contkey`;
- interface `ContentKeyable`;
- domain fields and methods use `CK` (for example, `ExpCK`, `SupExpCK`, `ValExpCK`, `GetExpCK`, and `SelectRecByCK`);
- local variables and structured log attributes use `ck` / `expCK`.

### PostgreSQL schema

Use `exp_ck` and `sup_exp_ck` as the canonical database column names:

- fresh databases get the new names from `000001_initial.up.sql`;
- existing databases are upgraded by `000002_content_key.up.sql`, which renames the old columns in both `proc_*` and `pool_*` tables;
- SQL queries and Go `db` tags use the new names.

The migration must be applied before deploying the engine version that expects the renamed columns. This is a coordinated schema/application change; the old application version is not expected to work against the renamed schema.

### JSON and SDK compatibility

Do not mechanically rename every serialized property:

- internal type-expression reference DTOs use `exp_ck`;
- the existing `termvar.VarRecDS` JSON property `vk` and term-expression properties `val_exp_vk` remain unchanged because they are existing serialized representations;
- the pinned `orglang/go-sdk` still exposes `ExpVK` on `proc/typeexp.ExpRef`, `pool/typeexp.ExpRef`, and `adt/termvar.VarRec`, with its existing JSON/query tags. The engine adapters explicitly translate between those SDK fields and internal `ExpCK` fields. The adjacent SDK repository is not changed by this engine refactor.

## Compatibility and migration

Existing PostgreSQL installations require migration `000002_content_key.up.sql`; fresh installations use the updated initial schema. No JSON migration is introduced for the legacy `vk` and `val_exp_vk` properties because those names remain at their existing serialization boundaries. SDK wire and query contracts remain unchanged through explicit adapter mappings.

## Consequences

- Go source code and the PostgreSQL schema use the canonical Content Key terminology.
- Existing database installations have an explicit column-rename migration.
- Existing serialized payloads and the pinned SDK contract retain their legacy names where they form compatibility boundaries.
- The legacy `vk`, `val_exp_vk`, and SDK `ExpVK` spellings are intentional boundary names, not canonical Go or database names.
