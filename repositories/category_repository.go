package repositories

import (
	"andre_kasir_api/models"
	"database/sql"
	"fmt"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query(`SELECT id, name, description FROM categories ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}

	return categories, rows.Err()
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow(
		`SELECT id, name, description FROM categories WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.Name, &c.Description)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &c, nil
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.QueryRow(
		`INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`,
		category.Name, category.Description,
	).Scan(&category.ID)
}

func (r *CategoryRepository) Update(category *models.Category) error {
	result, err := r.db.Exec(
		`UPDATE categories SET name = $1, description = $2 WHERE id = $3`,
		category.Name, category.Description, category.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *CategoryRepository) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}
