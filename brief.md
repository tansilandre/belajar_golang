# Brief: Kasir API Refactor

## What We're Doing
Transforming a single-file in-memory Go API into a production-ready application with proper architecture, persistent database, and new features (search, transactions, reporting).

## Context Before Starting

### Current State
- **Single file**: All logic lives in `main.go` (~260 lines)
- **No database**: Data stored in Go slices — resets on every restart
- **No config**: Port (`5001`) hardcoded, no `.env`
- **No dependencies**: Pure Go standard library
- **No structure**: No separation of concerns — handlers, models, logic all in one file
- **Module name**: `andre_kasir_api`

### Existing Endpoints
| Method | Route | Description |
|--------|-------|-------------|
| GET/POST | `/api/produk` | List / Create product |
| GET/PUT/DELETE | `/api/produk/{id}` | Read / Update / Delete product |
| GET/POST | `/api/categories` | List / Create category |
| GET/PUT/DELETE | `/api/categories/{id}` | Read / Update / Delete category |
| GET | `/health` | Health check |
| GET | `/` | Redirect to `/health` |

### Existing Models (in-memory)
- **Produk**: ID, Nama, Harga, Stok
- **Category**: ID, Name, Description

### Deployment
- VPS via GitHub Actions (push to `main` → SSH → build → systemd restart)
- Systemd service running as `ubuntu` user

## Key Decisions
- **PostgreSQL directly** instead of Supabase (as specified by user)
- **Viper** for config management (`.env` file)
- **`github.com/lib/pq`** as Postgres driver (v1.10.9)
- **Layered architecture**: handlers → services → repositories → database
- **Standard library `net/http`** for routing (no framework change)

## What Changes After
- Proper folder structure with clear separation of concerns
- Persistent data in PostgreSQL
- Config via `.env` (PORT, DB_CONN)
- Product search by name
- Checkout/transaction system with stock deduction
- Daily sales report endpoint
