package models

import (
	"frappuchino/internal/apperrors"
)

type MenuItem struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Allergens   []string `json:"allergens"`
	Size        string   `json:"size"`
}

type MenuItemIngredient struct {
	ID           int     `json:"id"`
	MenuItemID   string  `json:"menu_item_id"`
	Quantity     float64 `json:"quantity"`
	IngredientID string  `json:"ingredient_id"`
}

func NewMenuItem(allergens []string, dto CreateMenuRequest) (*MenuItem, error) {
	if dto.Name == "" || dto.Price <= 0 || dto.Size == "" {
		return nil, apperrors.ErrInvalidInput
	}

	description := dto.Description
	if description == "" {
		description = "No description"
	}

	return &MenuItem{
		ID:          dto.ID,
		Name:        dto.Name,
		Description: description,
		Price:       dto.Price,
		Allergens:   allergens,
		Size:        dto.Size,
	}, nil
}

func NewMenuItemIngredients(menuItemID string, items []MenuItemIngredientInput) ([]*MenuItemIngredient, error) {
	if menuItemID == "" || len(items) < 1 {
		return nil, apperrors.ErrInvalidInput
	}

	menuItems := []*MenuItemIngredient{}
	for _, item := range items {
		if item.Quantity <= 0 || item.IngredientID == "" {
			return nil, apperrors.ErrInvalidInput
		}
		menuItems = append(menuItems, &MenuItemIngredient{
			MenuItemID:   menuItemID,
			Quantity:     item.Quantity,
			IngredientID: item.IngredientID,
		})
	}
	return menuItems, nil
}
