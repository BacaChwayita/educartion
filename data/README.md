# Data Layer

## PostgreSQL setup scripts

SQL scripts are organized under `data\scripts` and should be run in order:

1. `data\scripts\database\001_create_database.sql`
2. `data\scripts\database\002_create_schemas.sql`
3. `data\scripts\schema\001_create_tables_app.sql`
4. `data\scripts\schema\002_create_tables_app_test.sql`

### Execution notes

- Run `001_create_database.sql` while connected to an admin database (commonly `postgres`).
- Then connect to `educartion_db` and run the schema scripts you need.
- Scripts are idempotent (`IF NOT EXISTS`) to support repeatable local and CI setup.
- `000_create_tables_for_schema.sql` is a shared SQL include used by the per-schema SQL scripts.

Example with local `psql`:

```powershell
psql -d postgres -f data\scripts\database\001_create_database.sql
psql -d educartion_db -f data\scripts\database\002_create_schemas.sql
psql -d educartion_db -f data\scripts\schema\001_create_tables_app.sql
psql -d educartion_db -f data\scripts\schema\002_create_tables_app_test.sql
```

## Docker scripts

Docker scripts are in `data\scripts\docker` and run PostgreSQL plus the SQL files for each schema.

1. App schema: `data\scripts\docker\run_app_schema.ps1`
2. Unit-test schema: `data\scripts\docker\run_app_test_schema.ps1`

They use this Compose file: `data\scripts\docker\docker-compose.postgres.yml`

Example:

```powershell
powershell -ExecutionPolicy Bypass -File data\scripts\docker\run_app_schema.ps1
powershell -ExecutionPolicy Bypass -File data\scripts\docker\run_app_test_schema.ps1
```
