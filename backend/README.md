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

## API Endpoints

Concise reference of implemented routes in `main.go`.

| Endpoint | Method | Auth | Input fields | Output fields | Success | Possible failure statuses |
| --- | --- | --- | --- | --- | --- | --- |
| `/api/auth/register` | POST | No | Body: `full_name`, `email`, `password_text` (`password_hash` JSON key in current model) | Empty body | `201 Created` | `400`, `500` |
| `/api/auth/login` | POST | No | Body: `email`, `password_text` (`password_hash` JSON key in current model) | `account_id`, `full_name`, `email`, `role`, `token` | `200 OK` | `400`, `401`, `404`, `409`, `500` |
| `/api/auth/logout` | POST | No | Body: `token` | Empty body | `204 No Content` | `400`, `404`, `500` |
| `/api/auth/me` | GET | No | Body currently decoded as `token` | No stable success payload yet | _No stable success contract yet_ | `400` |
| `/api/products` | GET | No | No body (handler currently attempts to decode optional search body) | Array of `product` objects | `200 OK` | `400`, `500` |
| `/api/products/get?id={id}` | GET | No | Query: `id` | One `product` object | `200 OK` | `400`, `404` |
| `/api/products/create` | POST | No | Body: `supplier_id`, `name`, `description`, `price`, `discount_percent`, `stock_quantity`, `is_active` | Created `product` object (includes `product_id`) | `201 Created` | `400`, `500` |
| `/api/products/update?id={id}` | PUT | No | Query: `id`; Body: product fields to update (`supplier_id`, `name`, `description`, `price`, `discount_percent`, `stock_quantity`, `is_active`) | `rows_affected` | `200 OK` | `400`, `500` |
| `/api/products/delete?id={id}` | DELETE | No | Query: `id` | `rows_affected` | `200 OK` | `400`, `500` |
| `/api/cart` | GET | Bearer token | Header: `Authorization: Bearer <token>` | `cart_id`, `account_id`, `session_key`, `total_price`, `created_at`, `updated_at` | `200 OK` | `401`, `404`, `405`, `500` |
| `/api/cart/items` | GET | Bearer token | Header: `Authorization: Bearer <token>` | Array of cart-item entries | `200 OK` | `401`, `404`, `405`, `500` |
| `/api/cart/items` | POST | Bearer token | Header: `Authorization: Bearer <token>`; Body: `product_id`, `quantity` | Created/updated cart item | `201 Created` | `400`, `401`, `404`, `405`, `500` |
| `/api/cart/item/{productId}` | PUT | Bearer token | Header: `Authorization: Bearer <token>`; Path: `productId`; Body: `quantity` | Updated cart item | `200 OK` | `400`, `401`, `404`, `405`, `500` |
| `/api/cart/item/{productId}` | DELETE | Bearer token | Header: `Authorization: Bearer <token>`; Path: `productId` | Empty body | `204 No Content` | `400`, `401`, `404`, `405`, `500` |
| `/api/orders/{id}` | GET | No | Path: `id` | `order_id`, `order_number`, `status`, `subtotal_amount`, `discount_amount`, `total_amount`, `placed_at`, `order_item[]` (with nested `product_details` and `supplier_details`) | `200 OK` | `400`, `404`, `500` |
