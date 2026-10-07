\set ON_ERROR_STOP on

SELECT format(
    'CREATE ROLE %I LOGIN PASSWORD %L',
    :'DB_OWNER',
    :'DB_OWNER_PASSWORD'
)
WHERE NOT EXISTS (
    SELECT 1 FROM pg_roles WHERE rolname = :'DB_OWNER'
) \gexec

SELECT format(
    'CREATE DATABASE %I OWNER %I',
    :'DB_NAME',
    :'DB_OWNER'
)
WHERE NOT EXISTS (
    SELECT 1 FROM pg_database WHERE datname = :'DB_NAME'
) \gexec

\connect :DB_NAME

SELECT format(
    'CREATE SCHEMA %I AUTHORIZATION %I',
    :'DB_SCHEMA',
    :'DB_OWNER'
)
WHERE NOT EXISTS (
    SELECT 1 FROM pg_namespace WHERE nspname = :'DB_SCHEMA'
) \gexec
