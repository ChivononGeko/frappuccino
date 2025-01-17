package models

import (
	"frappuchino/internal/apperrors"
)

type CreateInventoryRequest struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	StockLevel float64 `json:"stock_level"`
	Price      float64 `json:"price"`
	UnitType   string  `json:"unit_type"`
}

func NewCreateInventoryRequest(inventoryRequest CreateInventoryRequest) (*CreateInventoryRequest, error) {
	if inventoryRequest.Name == "" || inventoryRequest.UnitType == "" || inventoryRequest.StockLevel <= 0 || inventoryRequest.Price <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	if inventoryRequest.ID == "" {
		inventoryRequest.ID = fromNameToID(inventoryRequest.Name)
	}

	return &CreateInventoryRequest{
		ID:         inventoryRequest.ID,
		Name:       inventoryRequest.Name,
		StockLevel: inventoryRequest.StockLevel,
		Price:      inventoryRequest.Price,
		UnitType:   inventoryRequest.UnitType,
	}, nil
}
