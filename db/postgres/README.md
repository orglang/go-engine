# PostgreSQL migrations

The PostgreSQL migration image contains `psql`, the golang-migrate CLI, the DBA bootstrap SQL, and the owner migration.

## Runtime parameters

The Compose workflow accepts:

- `DBA_DSN` — administrative PostgreSQL DSN.
- `DB_NAME` — target database name.
- `DB_SCHEMA` — target schema name.
- `DB_OWNER` — owner role name.
- `DB_OWNER_PASSWORD` — password used only when the owner role is created.
- `OWNER_DSN` — DSN for the owner migration stage.

The defaults match the local development stack: database `orglang`, schema `orglang`, owner `orglang`, owner password `orglang`, and PostgreSQL admin credentials `postgres` / `password`. Credentials and identifiers can be overridden through the environment. URL credentials must be URL-encoded when they contain reserved characters.

`OWNER_DSN` uses `search_path=orglang` and `x-multi-statement=true` by default. For a different schema, set `OWNER_DSN` explicitly with `search_path=<DB_SCHEMA>`. This keeps golang-migrate's `schema_migrations` table in the target schema and executes the multi-statement migration transactionally.

## Build

From the `go-engine` repository root:

    docker build -f db/postgres/migrations.Dockerfile -t orglang/pg-operator:latest db/postgres

## Run

The existing database task remains the entry point:

    task db:process

The Compose workflow starts PostgreSQL, then the `database` one-shot service runs `dba/bootstrap.sql` with the administrative DSN. After it exits successfully, the `schema` one-shot service runs golang-migrate with the filesystem source `file:///migrations/owner` and the owner DSN.

The same migration image is used for both stages; only the command and credentials differ. Repeated runs are safe: bootstrap only creates missing objects, while golang-migrate records the applied migration in `schema_migrations` and does not reapply it.