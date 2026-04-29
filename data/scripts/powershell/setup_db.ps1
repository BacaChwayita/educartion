Write-Host "Starting DB Setup Scripts..."

# Create DB + Schemas
Write-Host "Create DB"
Get-Content .\scripts\database\001_create_database.sql | docker exec -i educartion-postgres psql -U postgres -d educartion_db
Write-Host "Create Schemas"
Get-Content .\scripts\database\002_create_schemas.sql | docker exec -i educartion-postgres psql -U postgres -d educartion_db

# Create Tables
Write-Host "Create Tables"
Get-Content .\scripts\schema\001_create_tables_app.sql | docker exec -i educartion-postgres psql -U postgres -d educartion_db
Write-Host "Create Test Tables"
Get-Content .\scripts\schema\002_create_tables_app_test.sql | docker exec -i educartion-postgres psql -U postgres -d educartion_db

# Check tables created
Write-Host "Check Tables"
docker exec -it educartion-postgres psql -U postgres -d educartion_db -c "\dt app.*"
Write-Host "Check Test Tables"
docker exec -it educartion-postgres psql -U postgres -d educartion_db -c "\dt app_test.*"

# Done
Write-Host "Done"
