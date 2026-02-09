package repositories

import (
	"andre_kasir_api/models"
	"database/sql"
	"fmt"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(searchName string) ([]models.Product, error) {
	var query string
	var args []interface{}

	if searchName != "" {
		query = `SELECT id, name, price, stock, category_id FROM products WHERE name ILIKE $1 ORDER BY id`
		args = append(args, "%"+searchName+"%")
	} else {
		query = `SELECT id, name, price, stock, category_id FROM products ORDER BY id`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, rows.Err()
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		`SELECT id, name, price, stock, category_id FROM products WHERE id = $1`,
		id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &p, nil
}

func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.QueryRow(
		`INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		product.Name, product.Price, product.Stock, product.CategoryID,
	).Scan(&product.ID)
}

func (r *ProductRepository) Update(product *models.Product) error {
	result, err := r.db.Exec(
		`UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5`,
		product.Name, product.Price, product.Stock, product.CategoryID, product.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
