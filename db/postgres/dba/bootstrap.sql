\set ON_ERROR_STOP on

SELECT format(
    'CREATE ROLE %I LOGIN PASSWORD %L',
    :'OWNER_USER',
    :'OWNER_PASSWORD'
)
WHERE NOT EXISTS (
    SELECT 1 FROM pg_roles WHERE rolname = :'OWNER_USER'
) \gexec

SELECT format(
    'CREATE DATABASE %I OWNER %I',
    :'DB_NAME',
    :'OWNER_USER'
)
WHERE NOT EXISTS (
    SELECT 1 FROM pg_database WHERE datname = :'DB_NAME'
) \gexec

\connect :DB_NAME

SELECT format(
    'CREATE SCHEMA %I AUTHORIZATION %I',
    :'DB_SCHEMA',
    :'OWNER_USER'
)
WHERE NOT EXISTS (
    SELECT 1 FROM pg_namespace WHERE nspname = :'DB_SCHEMA'
) \gexec
