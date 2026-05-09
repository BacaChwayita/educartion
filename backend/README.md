# Backend Layer

## Overview

The backend is a Go service that loads environment-based configuration, connects to PostgreSQL, and exposes HTTP endpoints.

## Backend file tree

```text
backend\
├── config\        # Loads app settings, DB config, and logger options from env vars.
├── model\         # Shared DTOs and service-level type definitions (ie Request+Response types).
├── rpc\           # HTTP/RPC utility helpers (for example JSON response helpers).
└── service\       # Core business logic and DB-access services.
    ├── auth\      # Authentication use-cases and auth-focused service logic.
    └── db\        # Database-facing CRUD/service operations for domain entities.
```

## Run locally

1. Copy `example.env` to `.env` and set values for your local database.
2. From `backend\`, run:

```powershell
go run .
```

Current API route scaffolded in `main.go`:

- `POST /api/auth/register`

