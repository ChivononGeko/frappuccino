package handler

import (
	"encoding/json"
	"frappuchino/internal/models"
	"log/slog"
	"net/http"
)

type InventoryService interface {
	CreateInventoryItemService(invent models.CreateInventoryRequest) error
	GetAllInventoryItemsService() ([]*models.InventoryItem, error)
	GetInventoryItemService(id string) (*models.InventoryItem, error)
	UpdateInventoryItemService(id string, inventoryItem models.CreateInventoryRequest) error
	DeleteInventoryItemService(id string) error
	GetLeftOversService(sortBy, page, pageSize string) (map[string]interface{}, error)
}

type InventoryHandler struct {
	inventoryService InventoryService
}

func NewInventHandler(iS InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: iS}
}

func (h *InventoryHandler) CreateInventoryItem(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		slog.Error("Data is not JSON format")
		return
	}

	var inputInvent models.CreateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&inputInvent); err != nil {
		slog.Error("Handler error in Create Inventory: decoding JSON data", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	invent, err := models.NewCreateInventoryRequest(inputInvent)
	if err != nil {
		slog.Error("Handler error in Create Inventory: invalid input data", "item", inputInvent, "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.inventoryService.CreateInventoryItemService(*invent); err != nil {
		status := mapAppErrorToStatus(err)
		slog.Error("Handler error in Create Inventory: creating inventory item", "inventory item", invent, "Error", err)
		writeError(w, err.Error(), status)
		return
	}

	slog.Info("Inventory created successfully", "inventory ID", invent.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *InventoryHandler) GetAllInventoryItems(w http.ResponseWriter, r *http.Request) {
	allInvents, err := h.inventoryService.GetAllInventoryItemsService()
	if err != nil {
		slog.Error("Handler error in Get Inventory: retrieving all inventory items", "error", err)
		writeError(w, "Failed to retrieve inventory items", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, allInvents)
	slog.Info("Inventory items retrieved successfully", "count", len(allInvents))
}

func (h *InventoryHandler) GetInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	inventId, err := h.inventoryService.GetInventoryItemService(id)
	if err != nil {
		status := mapAppErrorToStatus(err)
		slog.Error("Handler error in Get Inventory: retrieving inventory item", "id", id, "error", err)
		writeError(w, err.Error(), status)
		return
	}

	writeJSON(w, http.StatusOK, inventId)
	slog.Info("Inventory item retrieved successfully", "id", id)
}

func (h *InventoryHandler) UpdateInventoryItem(w http.ResponseWriter, r *http.Request) {
	if !isJSONFile(w, r) {
		slog.Error("Data is not JSON format")
		return
	}
	id := r.PathValue("id")

	var inputInvent models.CreateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&inputInvent); err != nil {
		slog.Error("Handler error in Update Inventory: decoding JSON data", "error", err)
		writeError(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	inventoryItem, err := models.NewCreateInventoryRequest(inputInvent)
	if err != nil {
		slog.Error("Handler error in Update Inventory: invalid input data", "item", inputInvent, "error", err)
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.inventoryService.UpdateInventoryItemService(id, *inventoryItem); err != nil {
		status := mapAppErrorToStatus(err)
		slog.Error("Handler error in Update Inventory: updating inventory", "inventory item", inventoryItem, "error", err)
		writeError(w, err.Error(), status)
		return
	}
	slog.Info("Inventory updated successfully", "id", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *InventoryHandler) DeleteInventoryItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.inventoryService.DeleteInventoryItemService(id); err != nil {
		status := mapAppErrorToStatus(err)
		slog.Error("Handler error in Delete Inventory: deleting inventory", "id", id, "error", err)
		writeError(w, err.Error(), status)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	slog.Info("Inventory item deleted successfully", "id", id)
}

func (h *InventoryHandler) GetLeftOvers(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sortBy")

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	pageSize := r.URL.Query().Get("pageSize")
	if pageSize == "" {
		pageSize = "10"
	}

	leftOvers, err := h.inventoryService.GetLeftOversService(sortBy, page, pageSize)
	if err != nil {
		slog.Error("Handler error in Get LeftOvers: retrieving left overs", "sortBy", sortBy, "page", page, "pageSize", pageSize, "error", err)
		writeError(w, "Failed to retrieve left overs", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, leftOvers)
	slog.Info("Left overs retrieved successfully")
}
