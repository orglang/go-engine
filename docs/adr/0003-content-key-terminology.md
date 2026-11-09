# ADR-0003: Content Key terminology and persistence compatibility

- **Status:** Accepted
- **Date:** 2026-11-09

## Context

The domain concept previously named Value Key is a content-derived identity. The old terminology is misleading, but the key's calculation, composition, inversion, and lifecycle are not changing.

The Go package and symbols are internal implementation details. PostgreSQL column names, JSON property names, and SDK fields are persisted or externally consumed representations and may already exist in deployed databases and clients.

## Decision

Use **Content Key** as the canonical domain terminology in Go:

- package `internal/adt/contkey`;
- interface `ContentKeyable`;
- fields and methods use `CK` (for example, `ExpCK`, `SupExpCK`, `ValExpCK`, `GetExpCK`, and `SelectRecByCK`);
- local variables and structured log attributes use `ck` / `expCK`.

Do not rename persisted or externally versioned identifiers as part of this terminology refactor:

- PostgreSQL columns and query identifiers such as `exp_vk` and `sup_exp_vk` remain unchanged;
- JSON property names `vk`, `exp_vk`, and `val_exp_vk` remain unchanged, even where the corresponding Go fields are now named `ExpCK` or `ValExpCK`;
- database and JSON struct tags continue to map the new Go names to the existing storage/wire names;
- the pinned `orglang/go-sdk` still exposes `ExpVK` on `proc/typeexp.ExpRef`, `pool/typeexp.ExpRef`, and `adt/termvar.VarRec`. The engine's type-expression adapters and term-decoding converter explicitly map internal `ExpCK` to and from those legacy SDK fields. The adjacent SDK repository is not changed by this internal engine refactor.

## Compatibility and migration

No database migration is required: the physical schema is intentionally unchanged, and existing rows remain readable and writable through the unchanged SQL identifiers. Existing JSON producers and consumers continue to exchange the same `vk`, `exp_vk`, and `val_exp_vk` properties. Renaming those identifiers later would be a separate compatibility change requiring an explicit migration/versioning plan and coordinated client rollout.

## Consequences

- Go source code uses one canonical domain term without changing key semantics.
- Existing databases and serialized payloads remain compatible.
- The legacy `vk`, `exp_vk`, and `val_exp_vk` spellings are intentionally limited to persistence and wire-format boundaries; they are not canonical names for new Go identifiers.
