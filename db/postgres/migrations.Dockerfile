FROM migrate/migrate:v4.19.1 AS migrate

FROM postgres:17.11-alpine3.24

COPY --from=migrate /migrate /usr/local/bin/migrate
COPY dba/bootstrap.sql /migrations/dba/bootstrap.sql
COPY owner /migrations/owner
