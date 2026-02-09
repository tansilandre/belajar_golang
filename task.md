# Task: Refactor & Extend Kasir API

## Current State
- Single-file Go API (`main.go`) with in-memory CRUD for products & categories
- No external dependencies, no database, no config management
- Module: `andre_kasir_api`

---

## Task 1: Layered Architecture + Database + Config (Session 2)

### 1.1 Project Structure
Restructure into layered architecture:
```
├── main.go              # entry point, dependency injection, routes
├── .env                 # config (PORT, DB_CONN)
├── database/
│   └── database.go      # PostgreSQL connection (InitDB)
├── models/
│   └── product.go       # Product struct
│   └── category.go      # Category struct
├── handlers/
│   └── product_handler.go
│   └── category_handler.go
├── services/
│   └── product_service.go
│   └── category_service.go
├── repositories/
│   └── product_repository.go
│   └── category_repository.go
```

### 1.2 Config with Viper
- Install `github.com/spf13/viper`
- Read `.env` file for `PORT` and `DB_CONN`
- Config struct: `Port string`, `DBConn string`

### 1.3 Database (PostgreSQL — NOT Supabase)
- Install `github.com/lib/pq`
- `database/database.go`: `InitDB(connStr) (*sql.DB, error)` with connection pool
- `.env` variables needed:
  ```
  PORT=5001
  DB_CONN=postgresql://user:password@host:port/dbname
  ```

### 1.4 Products — Layered (model, handler, service, repo)
- Model: `Product { ID int, Name string, Price int, Stock int }`
- Repository: SQL queries against `products` table (GetAll, GetByID, Create, Update, Delete)
- Service: passes through to repo
- Handler: HTTP method routing, JSON encode/decode
- Routes: `GET/POST /api/produk`, `GET/PUT/DELETE /api/produk/{id}`

### 1.5 Categories — Layered (THE ACTUAL TASK)
- Same pattern as products for `categories` table
- Model: `Category { ID int, Name string, Description string }`
- Full CRUD: GetAll, GetByID, Create, Update, Delete
- Routes: `GET/POST /api/categories`, `GET/PUT/DELETE /api/categories/{id}`

### 1.6 Dependency Injection in main.go
```
repo → service → handler (for each resource)
```

---

## Task 2: Search + Transactions + Report (Session 3)

### 2.1 Search Products by Name
- `GET /api/produk?name=indom`
- Handler: read `name` from query params
- Repository: add `WHERE p.name ILIKE $1` with `%name%` when name is provided

### 2.2 Transaction Tables (SQL to run manually or via migration)
```sql
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    total_amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS transaction_details (
    id SERIAL PRIMARY KEY,
    transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    subtotal INT NOT NULL
);
```

### 2.3 Transaction Models
- `Transaction { ID, TotalAmount, CreatedAt, Details []TransactionDetail }`
- `TransactionDetail { ID, TransactionID, ProductID, ProductName, Quantity, Subtotal }`
- `CheckoutItem { ProductID, Quantity }`
- `CheckoutRequest { Items []CheckoutItem }`

### 2.4 Checkout API
- `POST /api/checkout` with body `{ "items": [{ "product_id": 1, "quantity": 2 }] }`
- Flow (inside DB transaction):
  1. For each item: fetch product price/stock, calculate subtotal
  2. Deduct stock (`UPDATE products SET stock = stock - qty`)
  3. Insert into `transactions` (total_amount) → get transaction ID
  4. Insert each detail into `transaction_details`
  5. Commit
- Layered: handler → service → repository

### 2.5 Sales Report API
- `GET /api/report/hari-ini`
- Response:
  ```json
  {
    "total_revenue": 45000,
    "total_transaksi": 5,
    "produk_terlaris": { "nama": "Indomie Goreng", "qty_terjual": 12 }
  }
  ```

### 2.6 Optional: Date Range Report
- `GET /api/report?start_date=2026-01-01&end_date=2026-02-01`

---

## .env Variables Needed
```
PORT=5001
DB_CONN=postgresql://user:password@host:port/dbname
```

## Dependencies to Install
```
go get github.com/spf13/viper
go get github.com/lib/pq@v1.10.9
```
