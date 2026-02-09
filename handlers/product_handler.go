package handlers

import (
	"andre_kasir_api/models"
	"andre_kasir_api/services"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getAll(w, r)
	case http.MethodPost:
		h.create(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *ProductHandler) HandleProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getByID(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (h *ProductHandler) getAll(w http.ResponseWriter, r *http.Request) {
	searchName := r.URL.Query().Get("name")

	products, err := h.service.GetAll(searchName)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) getByID(w http.ResponseWriter, r *http.Request, id int) {
	product, err := h.service.GetByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if product == nil {
		writeError(w, http.StatusNotFound, "Product not found")
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.Create(&product); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) update(w http.ResponseWriter, r *http.Request, id int) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	product.ID = id
	if err := h.service.Update(&product); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) delete(w http.ResponseWriter, r *http.Request, id int) {
	if err := h.service.Delete(id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeError(w, http.StatusNotFound, "Product not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]interface{}{
		"error": message,
		"code":  status,
	})
}
