# PostgreSQL migrations

The PostgreSQL migration image contains `psql`, the golang-migrate CLI, the DBA bootstrap SQL, and the owner migration.

## Runtime parameters

The Compose workflow accepts:

- `DBA_USER` — administrative PostgreSQL user.
- `DBA_PASSWORD` — administrative PostgreSQL password.
- `DBA_URL` — administrative PostgreSQL connection URL.
- `DB_NAME` — target database name.
- `DB_SCHEMA` — target schema name.
- `OWNER_USER` — owner role name.
- `OWNER_PASSWORD` — password used only when the owner role is created.
- `OWNER_URL` — connection URL for the owner migration stage.

These parameters have no defaults in the Compose file and must be provided by the environment. The local development stack defines them in `stack/ops/commons/.env`. URL credentials must be URL-encoded when they contain reserved characters.

`OWNER_URL` should use `search_path=<DB_SCHEMA>` and `x-multi-statement=true`. This keeps golang-migrate's `schema_migrations` table in the target schema and executes the multi-statement migration transactionally.

## Build

From the `go-engine` repository root:

    docker build -f db/postgres/migrations.Dockerfile -t orglang/pg-operator:latest db/postgres

## Run

The existing database task remains the entry point:

    task db:process

The Compose workflow starts PostgreSQL, then the `database` one-shot service runs `dba/bootstrap.sql` with the administrative URL. After it exits successfully, the `schema` one-shot service runs golang-migrate with the filesystem source `file:///migrations/owner` and the owner URL.

The same migration image is used for both stages; only the command and credentials differ. Repeated runs are safe: bootstrap only creates missing objects, while golang-migrate records the applied migration in `schema_migrations` and does not reapply it.

<!-- CI trigger -->
