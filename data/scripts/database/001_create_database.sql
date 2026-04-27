-- Run this script from an admin database connection (for example: postgres).
-- This uses psql's \gexec because CREATE DATABASE cannot run in a transaction block.
SELECT 'CREATE DATABASE educartion_db'
WHERE NOT EXISTS (
    SELECT 1
    FROM pg_database
    WHERE datname = 'educartion_db'
)\gexec;
