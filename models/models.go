package models

import "time"

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID *int   `json:"category_id,omitempty"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details,omitempty"`
}

type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

type SalesReport struct {
	TotalRevenue   int             `json:"total_revenue"`
	TotalTransaksi int             `json:"total_transaksi"`
	ProdukTerlaris *ProdukTerlaris `json:"produk_terlaris,omitempty"`
}

type ProdukTerlaris struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}
