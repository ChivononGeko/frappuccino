package models

import (
	"frappuchino/internal/apperrors"
)

type CreateMenuRequest struct {
	ID          string                    `json:"product_id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Price       float64                   `json:"price"`
	Size        string                    `json:"size"`
	Ingredients []MenuItemIngredientInput `json:"ingredients"`
}

type MenuItemIngredientInput struct {
	IngredientID string  `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
}

func NewCreateMenuRequest(menuRequest CreateMenuRequest) (*CreateMenuRequest, error) {
	if menuRequest.Name == "" || menuRequest.Price <= 0 || menuRequest.Size == "" {
		return nil, apperrors.ErrInvalidInput
	}

	if menuRequest.ID == "" {
		menuRequest.ID = fromNameToID(menuRequest.Name)
	}
	if menuRequest.Description == "" {
		menuRequest.Description = "No description"
	}

	for _, ingredient := range menuRequest.Ingredients {
		if ingredient.IngredientID == "" || ingredient.Quantity <= 0 {
			return nil, apperrors.ErrInvalidInput
		}
	}

	return &CreateMenuRequest{
		ID:          menuRequest.ID,
		Name:        menuRequest.Name,
		Description: menuRequest.Description,
		Price:       menuRequest.Price,
		Size:        menuRequest.Size,
		Ingredients: menuRequest.Ingredients,
	}, nil
}
