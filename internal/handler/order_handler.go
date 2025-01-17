package handler

import (
	"encoding/json"
	"frappuchino/internal/models"
	"log/slog"
	"net/http"
	"strconv"
)

type OrderService interface {
	CreateOrderService(newOrder models.CreateOrderRequest) error
	GetAllOrdersService() ([]*models.Order, error)
	GetOrderService(id int) (*models.Order, error)
	UpdateOrderService(id int, updateOrder models.CreateOrderRequest) error
	DeleteOrderService(id int) error
	CloseOrderService(id int) error
	NumberOfOrderedItemsService(start, end string) (map[string]int, error)
	AddOrdersService(orders []models.CreateOrderRequest) error
}

type OrderHandler struct {
	orderService OrderService
}

func NewOrderHandler(os OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: os,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		slog.Error("Data is not JSON format")
		return
	}

	var inputOrder models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&inputOrder); err != nil {
		slog.Error("Handler error in Create Order: decoding JSON data", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	order, err := models.NewCreateOrder(inputOrder)
	if err != nil {
		slog.Error("Handler error in Create Order: invalid input data", "item", inputOrder, "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = h.orderService.CreateOrderService(*order); err != nil {
		status := mapAppErrorToStatus(err)
		slog.Error("Handler error in Create Order: creating order", "order", order, "error", err)
		writeError(w, err.Error(), status)
		return
	}

	slog.Info("Order created successfully", "Customer", order.CustomerName)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	allOrders, err := h.orderService.GetAllOrdersService()
	if err != nil {
		slog.Error("Handler error in Get Orders: retrieving all orders", "error", err)
		writeError(w, "Failed to retrieve all orders", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, allOrders)
	slog.Info("All orders retrieved successfully", "count", len(allOrders))
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("Handler error in Get Order: id type conversion", "id", r.PathValue("id"), "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderId, err := h.orderService.GetOrderService(id)
	if err != nil {
		status := mapAppErrorToStatus(err)
		slog.Error("Handler error in Get Order: retrieving order", "id", id, "error", err)
		writeError(w, err.Error(), status)
		return
	}
	writeJSON(w, http.StatusOK, orderId)
	slog.Info("Order retrieved successfully", "id", id)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		slog.Error("Data is not JSON format")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("Handler error in Update Order: id type conversion", "id", r.PathValue("id"), "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var inputOrder models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&inputOrder); err != nil {
		slog.Error("Handler error in Update Order: decoding JSON data", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	order, err := models.NewCreateOrder(inputOrder)
	if err != nil {
		slog.Error("Handler error in Update Order: invalid input data", "item", inputOrder, "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.orderService.UpdateOrderService(id, *order); err != nil {
		status := mapAppErrorToStatus(err)
		slog.Error("Handler error in Update Order: updating order", "order", order, "error", err)
		writeError(w, err.Error(), status)
		return
	}

	slog.Info("Order updated successfully", "id", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("Handler error in Delete Order: id type conversion", "id", r.PathValue("id"), "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.orderService.DeleteOrderService(id); err != nil {
		status := mapAppErrorToStatus(err)
		slog.Error("Handler error in Delete Order: deleting order", "id", id, "error", err)
		writeError(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	slog.Info("Order deleted successfully", "id", id)
}

func (h *OrderHandler) CloseOrder(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("Handler error in Close Order: id type conversion", "id", r.PathValue("id"), "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.orderService.CloseOrderService(id); err != nil {
		status := mapAppErrorToStatus(err)
		slog.Error("Handler error in Close Order: closing order", "id", id, "error", err)
		writeError(w, err.Error(), status)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Order closed successfully", "id", id)
}

func (h *OrderHandler) NumberOfOrderedItems(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	startDate := queryParams.Get("startDate")
	endDate := queryParams.Get("endDate")

	orderedItems, err := h.orderService.NumberOfOrderedItemsService(startDate, endDate)
	if err != nil {
		slog.Error("Handler error in Number Of Ordered Items: ", "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(orderedItems); err != nil {
		slog.Error("Handler error in Get Order: encoding JSON data", "error", err)
		writeError(w, "Failed to encode order to JSON", http.StatusInternalServerError)
		return
	}

	slog.Info("Number of Ordered items retrieved successfully")
}

func (h *OrderHandler) BatchCreateOrders(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		slog.Error("Data is not JSON format")
		writeError(w, "Invalid format", http.StatusBadRequest)
		return
	}

	var inputOrders []models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&inputOrders); err != nil {
		slog.Error("Handler error in Create Orders: decoding JSON data", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if err := h.orderService.AddOrdersService(inputOrders); err != nil {
		status := mapAppErrorToStatus(err)
		slog.Error("Handler error in Batch Create Orders: creating orders", "error", err)
		writeError(w, err.Error(), status)
		return
	}

	slog.Info("Orders created successfully", "orders_count", len(inputOrders))
	w.WriteHeader(http.StatusCreated)
}
