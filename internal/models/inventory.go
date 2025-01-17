package models

import (
	"frappuchino/internal/apperrors"
	"time"
)

type InventoryItem struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	StockLevel  float64   `json:"stock_level"`
	UnitType    string    `json:"unit_type"`
	Price       float64   `json:"price"`
	LastUpdated time.Time `json:"last_update"`
}

type InventoryTransaction struct {
	ID              int       `json:"id"`
	InventoryID     string    `json:"inventory_id"`
	ChangeAmount    float64   `json:"change_amount"`
	TransactionType string    `json:"transaction_type"`
	ChangeAt        time.Time `json:"occurred_at"`
}

func NewInventoryItem(dto CreateInventoryRequest) (*InventoryItem, error) {
	if dto.Name == "" || dto.StockLevel < 0 || dto.UnitType == "" || dto.Price <= 0 {
		return nil, apperrors.ErrInvalidInput
	}

	if dto.ID == "" {
		dto.ID = fromNameToID(dto.Name)
	}

	return &InventoryItem{
		ID:          dto.ID,
		Name:        dto.Name,
		StockLevel:  dto.StockLevel,
		UnitType:    dto.UnitType,
		Price:       dto.Price,
		LastUpdated: time.Now(),
	}, nil
}

func NewInventoryTransaction(inventoryID string, changeAmount float64, transactionType string) (*InventoryTransaction, error) {
	if inventoryID == "" || changeAmount == 0 {
		return nil, apperrors.ErrInvalidInput
	}

	if !(transactionType == "added" || transactionType == "written off" || transactionType == "sale" || transactionType == "created") {
		return nil, apperrors.ErrInvalidInput
	}

	if changeAmount < 0 && transactionType == "added" {
		transactionType = "written off"
	}

	if transactionType == "sale" {
		changeAmount *= (-1)
	}

	return &InventoryTransaction{
		InventoryID:     inventoryID,
		ChangeAmount:    changeAmount,
		TransactionType: transactionType,
		ChangeAt:        time.Now(),
	}, nil
}
