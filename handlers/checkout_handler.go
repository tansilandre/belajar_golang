package handlers

import (
	"andre_kasir_api/models"
	"andre_kasir_api/services"
	"encoding/json"
	"net/http"
)

type CheckoutHandler struct {
	service *services.TransactionService
}

func NewCheckoutHandler(service *services.TransactionService) *CheckoutHandler {
	return &CheckoutHandler{service: service}
}

func (h *CheckoutHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.Items) == 0 {
		writeError(w, http.StatusBadRequest, "Items cannot be empty")
		return
	}

	transaction, err := h.service.Checkout(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, transaction)
}
