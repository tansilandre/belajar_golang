package handlers

import (
	"andre_kasir_api/services"
	"net/http"
	"time"
)

type ReportHandler struct {
	service *services.TransactionService
}

func NewReportHandler(service *services.TransactionService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	path := r.URL.Path

	if path == "/api/report/hari-ini" {
		h.getDailyReport(w, r)
		return
	}

	if path == "/api/report" {
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")

		if startDateStr == "" || endDateStr == "" {
			writeError(w, http.StatusBadRequest, "start_date and end_date are required")
			return
		}

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid start_date format (YYYY-MM-DD)")
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid end_date format (YYYY-MM-DD)")
			return
		}

		report, err := h.service.GetReportByDateRange(startDate, endDate)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, report)
		return
	}

	writeError(w, http.StatusNotFound, "Report endpoint not found")
}

func (h *ReportHandler) getDailyReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetDailyReport()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, report)
}
