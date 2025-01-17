package handler

import (
	"frappuchino/internal/models"
	"log/slog"
	"net/http"
)

type ReportsService interface {
	TotalSalesReportService() (*models.TotalPrice, error)
	PopularItemsReportService() ([]*models.PopularItem, error)
	SearchService(q, filter, minPrice, maxPrice string) (map[string]interface{}, error)
	OrderedItemsByPeriodService(period, month, year string) (map[string]interface{}, error)
}

type ReportsHandler struct {
	reportsService ReportsService
}

func NewReportsHandler(rs ReportsService) *ReportsHandler {
	return &ReportsHandler{reportsService: rs}
}

func (h *ReportsHandler) TotalSalesReportHandler(w http.ResponseWriter, r *http.Request) {
	totalSales, err := h.reportsService.TotalSalesReportService()
	if err != nil {
		slog.Error("Handler error in Total Sales Report: counting sales", "error", err)
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Get total sales successful", "total sales", totalSales.TotalSale)
	writeJSON(w, http.StatusOK, totalSales)
}

func (h *ReportsHandler) PopularItemsReportHandler(w http.ResponseWriter, r *http.Request) {
	popularItems, err := h.reportsService.PopularItemsReportService()
	if err != nil {
		slog.Error("Handler error in Popular Items Report: identifying items", "error", err)
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Get popular items successful")
	writeJSON(w, http.StatusOK, map[string][]*models.PopularItem{"the most popular item:": popularItems})
}

func (h *ReportsHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	q := queryParams.Get("q")
	if q == "" {
		slog.Error("Handler error in Search: missing required parameter q")
		writeError(w, "Missing required parameter: q", http.StatusBadRequest)
		return
	}

	filter := queryParams.Get("filter")
	if filter == "" || filter == "menu,orders" || filter == "orders,menu" {
		filter = "all"
	}

	minPrice := queryParams.Get("minPrice")
	if minPrice == "" {
		minPrice = "0"
	}
	maxPrice := queryParams.Get("maxPrice")
	if maxPrice == "" {
		maxPrice = "1000000"
	}

	response, err := h.reportsService.SearchService(q, filter, minPrice, maxPrice)
	if err != nil {
		slog.Error("Handler error in Search: failed retrieved items", "q", q, "filter", filter, "min price", minPrice, "max price", maxPrice, "error", err)
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Search items successful")
	writeJSON(w, http.StatusOK, response)
}

func (h *ReportsHandler) OrderedItemsByPeriodHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	period := queryParams.Get("period")
	if period == "" {
		slog.Error("Handler error from Ordered Items by Period: missing required parameter period")
		writeError(w, "Missing required parameter period", http.StatusBadRequest)
		return
	}

	if !(period == "day" || period == "month") {
		slog.Error("Handler error from Ordered Items by Period: invalid value for required parameter period", "period", period)
		writeError(w, "Invalid value for required parameter period", http.StatusBadRequest)
		return
	}

	month := queryParams.Get("month")
	year := queryParams.Get("year")

	response, err := h.reportsService.OrderedItemsByPeriodService(period, month, year)
	if err != nil {
		slog.Error("Handler error in Ordered Items by Period: failed retrieved items", "period", period, "month", month, "year", year, "error", err)
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Get ordered items by period successful")
	writeJSON(w, http.StatusOK, response)
}
