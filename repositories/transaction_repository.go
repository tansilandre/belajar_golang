package repositories

import (
	"andre_kasir_api/models"
	"database/sql"
	"fmt"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Checkout(req *models.CheckoutRequest) (*models.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var totalAmount int
	details := make([]models.TransactionDetail, 0, len(req.Items))

	for _, item := range req.Items {
		var price, stock int
		var productName string
		err := tx.QueryRow(
			`SELECT name, price, stock FROM products WHERE id = $1 FOR UPDATE`,
			item.ProductID,
		).Scan(&productName, &price, &stock)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product with ID %d not found", item.ProductID)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get product %d: %w", item.ProductID, err)
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s: available %d, requested %d", productName, stock, item.Quantity)
		}

		subtotal := price * item.Quantity
		totalAmount += subtotal

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})

		_, err = tx.Exec(
			`UPDATE products SET stock = stock - $1 WHERE id = $2`,
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock for product %d: %w", item.ProductID, err)
		}
	}

	var transactionID int
	err = tx.QueryRow(
		`INSERT INTO transactions (total_amount, created_at) VALUES ($1, $2) RETURNING id`,
		totalAmount, time.Now(),
	).Scan(&transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	for i := range details {
		_, err = tx.Exec(
			`INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)`,
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create transaction detail: %w", err)
		}
		details[i].TransactionID = transactionID
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   time.Now(),
		Details:     details,
	}, nil
}

func (r *TransactionRepository) GetDailyReport(date time.Time) (*models.SalesReport, error) {
	report := &models.SalesReport{}

	err := r.db.QueryRow(
		`SELECT COALESCE(SUM(total_amount), 0), COUNT(*) FROM transactions WHERE DATE(created_at) = DATE($1)`,
		date,
	).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily report: %w", err)
	}

	var productName string
	var qtyTerjual int
	err = r.db.QueryRow(
		`SELECT p.name, SUM(td.quantity) as qty_terjual
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) = DATE($1)
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1`,
		date,
	).Scan(&productName, &qtyTerjual)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get top product: %w", err)
	}

	if err == nil {
		report.ProdukTerlaris = &models.ProdukTerlaris{
			Nama:       productName,
			QtyTerjual: qtyTerjual,
		}
	}

	return report, nil
}

func (r *TransactionRepository) GetReportByDateRange(startDate, endDate time.Time) (*models.SalesReport, error) {
	report := &models.SalesReport{}

	err := r.db.QueryRow(
		`SELECT COALESCE(SUM(total_amount), 0), COUNT(*) FROM transactions WHERE DATE(created_at) BETWEEN DATE($1) AND DATE($2)`,
		startDate, endDate,
	).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}

	var productName string
	var qtyTerjual int
	err = r.db.QueryRow(
		`SELECT p.name, SUM(td.quantity) as qty_terjual
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) BETWEEN DATE($1) AND DATE($2)
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1`,
		startDate, endDate,
	).Scan(&productName, &qtyTerjual)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get top product: %w", err)
	}

	if err == nil {
		report.ProdukTerlaris = &models.ProdukTerlaris{
			Nama:       productName,
			QtyTerjual: qtyTerjual,
		}
	}

	return report, nil
}
