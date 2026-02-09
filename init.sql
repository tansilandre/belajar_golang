-- Run this SQL as postgres superuser

CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    category_id INT REFERENCES categories(id) ON DELETE SET NULL
);

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

-- Grant permissions to asisten_intern user
GRANT ALL PRIVILEGES ON TABLE categories TO asisten_intern;
GRANT ALL PRIVILEGES ON TABLE products TO asisten_intern;
GRANT ALL PRIVILEGES ON TABLE transactions TO asisten_intern;
GRANT ALL PRIVILEGES ON TABLE transaction_details TO asisten_intern;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO asisten_intern;
