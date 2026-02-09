package main

import (
	"andre_kasir_api/config"
	"andre_kasir_api/database"
	"andre_kasir_api/handlers"
	"andre_kasir_api/repositories"
	"andre_kasir_api/services"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	db, err := database.InitDB(cfg.DBConn)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		return
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo)

	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	checkoutHandler := handlers.NewCheckoutHandler(transactionService)
	reportHandler := handlers.NewReportHandler(transactionService)

	http.HandleFunc("/api/produk/", productHandler.HandleProduct)
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategory)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/checkout", checkoutHandler.HandleCheckout)
	http.HandleFunc("/api/report/hari-ini", reportHandler.HandleReport)
	http.HandleFunc("/api/report", reportHandler.HandleReport)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "200",
			"message": "API running",
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/health", http.StatusMovedPermanently)
	})

	fmt.Printf("Server started on :%s\n", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
