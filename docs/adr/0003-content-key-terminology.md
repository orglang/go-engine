# ADR-0003: Content Key terminology and persistence compatibility

- **Status:** Accepted
- **Date:** 2026-11-09

## Context

The domain concept previously named Value Key is a content-derived identity. The old terminology is misleading, but the key's calculation, composition, inversion, and lifecycle are not changing.

The Go package and symbols are internal implementation details. The PostgreSQL column names and the JSON property name are persisted or externally consumed representations and may already exist in deployed databases and clients.

## Decision

Use **Content Key** as the canonical domain terminology in Go:

- package `internal/adt/contkey`;
- interface `ContentKeyable`;
- fields and methods use `CK` (for example, `ExpCK`, `SupExpCK`, `GetExpCK`, and `SelectRecByCK`);
- local variables and structured log attributes use `ck` / `expCK`.

Do not rename persisted identifiers as part of this terminology refactor:

- PostgreSQL columns and query identifiers such as `exp_vk`, `sup_exp_vk`, and `val_exp_vk` remain unchanged;
- the JSON property `vk` remains unchanged, even though its Go field is now named `ExpCK`;
- database and JSON struct tags continue to map the new Go names to the existing storage/wire names.
- The pinned `orglang/go-sdk` contract still exposes `ExpVK` on `proc/typeexp.ExpRef` and `pool/typeexp.ExpRef`; adapters deliberately map internal `ExpCK` to and from that legacy SDK field. The adjacent SDK repository is not changed by this internal engine refactor.

## Compatibility and migration

No database migration is required: the physical schema is intentionally unchanged, and existing rows remain readable and writable through the unchanged SQL identifiers. Existing JSON producers and consumers continue to exchange the same `vk` property. Renaming those identifiers later would be a separate compatibility change requiring an explicit migration/versioning plan and coordinated client rollout.

## Consequences

- Go source code uses one canonical domain term without changing key semantics.
- Existing databases and serialized payloads remain compatible.
- The legacy `vk`/ `exp_vk` spellings are intentionally limited to persistence and wire-format boundaries; they are not canonical names for new Go identifiers.
