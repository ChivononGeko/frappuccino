package service

import (
	"frappuchino/internal/models"
	"log/slog"
	"strconv"
)

type InventoryRepository interface {
	AddInventoryItemRepository(inventoryItem models.InventoryItem, inventoryTransaction models.InventoryTransaction) error
	GetInventoryItemRepository(id string) (*models.InventoryItem, error)
	GetAllInventoryItemsRepository() ([]*models.InventoryItem, error)
	UpdateInventoryItemRepository(id string, inventoryItem models.InventoryItem, inventoryTransaction models.InventoryTransaction) error
	DeleteInventoryItemRepository(id string) error
	GetLeftOversRepository(sortBy string, page, offset, pageSize int) (map[string]interface{}, error)
}

type InventoryService struct {
	inventoryRepo InventoryRepository
}

func NewInventoryService(iR InventoryRepository) *InventoryService {
	return &InventoryService{inventoryRepo: iR}
}

func (s *InventoryService) CreateInventoryItemService(inventoryItemRequest models.CreateInventoryRequest) error {
	inventoryItem, inventoryTransaction, err := s.createInventoryObjects(inventoryItemRequest, "created")
	if err != nil {
		slog.Error("Service error in Create Inventory: failed to create objects", "error", err)
		return err
	}

	err = s.inventoryRepo.AddInventoryItemRepository(*inventoryItem, *inventoryTransaction)
	if err != nil {
		slog.Error("Service error in Create Inventory: failed to add data to tables", "item", inventoryItem, "transaction", inventoryTransaction, "error", err)
		return err
	}

	return nil
}

func (s *InventoryService) GetAllInventoryItemsService() ([]*models.InventoryItem, error) {
	inventoryItems, err := s.inventoryRepo.GetAllInventoryItemsRepository()
	if err != nil {
		slog.Error("Service error in Get Inventory: failed to retrieve all inventory items", "error", err)
		return nil, err
	}
	return inventoryItems, nil
}

func (s *InventoryService) GetInventoryItemService(id string) (*models.InventoryItem, error) {
	inventoryItem, err := s.inventoryRepo.GetInventoryItemRepository(id)
	if err != nil {
		slog.Error("Service error in Get Inventory: failed to retrieve all inventory item", "id", id, "error", err)
		return nil, err
	}
	return inventoryItem, nil
}

func (s *InventoryService) UpdateInventoryItemService(id string, inventoryItemRequest models.CreateInventoryRequest) error {
	inventoryItemRequest.ID = id
	inventoryItem, inventoryTransaction, err := s.createInventoryObjects(inventoryItemRequest, "added")
	if err != nil {
		slog.Error("Service error in Update Inventory: failed to create objects", "error", err)
		return err
	}

	err = s.inventoryRepo.UpdateInventoryItemRepository(id, *inventoryItem, *inventoryTransaction)
	if err != nil {
		slog.Error("Service error in Update Inventory: failed to update inventory", "id", id, "error", err)
		return err
	}

	return nil
}

func (s *InventoryService) DeleteInventoryItemService(id string) error {
	err := s.inventoryRepo.DeleteInventoryItemRepository(id)
	if err != nil {
		slog.Error("Service error in Delete Inventory: failed to delete inventory", "id", id, "error", err)
		return err
	}
	return nil
}

func (s *InventoryService) createInventoryObjects(inventoryItemRequest models.CreateInventoryRequest, typeTransaction string) (*models.InventoryItem, *models.InventoryTransaction, error) {
	inventoryItem, err := models.NewInventoryItem(inventoryItemRequest)
	if err != nil {
		slog.Error("Service error in Create Object: failed to create inventory item", "input item", inventoryItemRequest, "error", err)
		return nil, nil, err
	}

	inventoryTransaction, err := models.NewInventoryTransaction(inventoryItemRequest.ID, inventoryItem.StockLevel, typeTransaction)
	if err != nil {
		slog.Error("Service error in Create Object: failed to create inventory transaction", "id", inventoryItemRequest.ID, "stock level", inventoryItemRequest.StockLevel, "type transaction", typeTransaction, "error", err)
		return nil, nil, err
	}

	return inventoryItem, inventoryTransaction, nil
}

func (s *InventoryService) GetLeftOversService(sortBy, pageParam, pageSizeParam string) (map[string]interface{}, error) {
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		slog.Error("Service error in Get Leftovers: failed to convert page to integer", "page", pageParam, "error", err)
		return nil, err
	}

	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil {
		slog.Error("Service error in Get Leftovers: failed to convert pageSize to integer", "page size", pageSizeParam, "error", err)
		return nil, err
	}

	offset := (page - 1) * pageSize

	leftovers, err := s.inventoryRepo.GetLeftOversRepository(sortBy, page, offset, pageSize)
	if err != nil {
		slog.Error("Service error in Get Leftovers: failed to retrieve leftovers", "sort by", sortBy, "offset", offset, "page size", pageSize, "error", err)
		return nil, err
	}

	return leftovers, nil
}
