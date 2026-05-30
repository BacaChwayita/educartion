#!/bin/bash
set -e

# Create test db. Main db made by default by postgres image.
createdb -U "$POSTGRES_USER" "$POSTGRES_DB_TEST"

# Create tables in both dbs.
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /sql/data/scripts/create_tables.sql
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB_TEST" -f /sql/data/scripts/create_tables.sql

# Create Dummy Data
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f /sql/data/scripts/seed/001_seed_data.sql
