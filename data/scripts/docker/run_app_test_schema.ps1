$ErrorActionPreference = "Stop"

$composeFile = "data\scripts\docker\docker-compose.postgres.yml"

docker compose -f $composeFile up -d postgres
docker compose -f $composeFile exec -T postgres psql -U postgres -d postgres -f /sql/database/001_create_database.sql
docker compose -f $composeFile exec -T postgres psql -U postgres -d educartion_db -f /sql/database/002_create_schemas.sql
docker compose -f $composeFile exec -T postgres psql -U postgres -d educartion_db -f /sql/schema/002_create_tables_app_test.sql

