package models

import (
	"encoding/json"
	"frappuchino/internal/apperrors"
	"time"
)

type Order struct {
	ID                  int             `json:"id"`
	CustomerID          int             `json:"customer_id"`
	TotalAmount         float64         `json:"total_amount"`
	Status              string          `json:"status"`
	SpecialInstructions json.RawMessage `json:"special_instructions"`
	PaymentMethod       string          `json:"payment_method"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
}

type OrderItem struct {
	ID         string  `json:"id"`
	OrderID    int     `json:"order_id"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price_at_order"`
	MenuItemID string  `json:"menu_item_id"`
}

func NewOrder(customerID int, totalAmount float64, dto CreateOrderRequest) (*Order, error) {
	if customerID < 1 || totalAmount == 0 || dto.PaymentMethod == "" {
		return nil, apperrors.ErrInvalidInput
	}

	if !(dto.PaymentMethod == "card" || dto.PaymentMethod == "cash" || dto.PaymentMethod == "kaspi_qr") {
		return nil, apperrors.ErrInvalidInput
	}

	if dto.Instructions == nil {
		dto.Instructions = json.RawMessage(`{"special_request":"no"}`)
	}
	return &Order{
		CustomerID:          customerID,
		TotalAmount:         totalAmount,
		Status:              "open",
		SpecialInstructions: dto.Instructions,
		PaymentMethod:       dto.PaymentMethod,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}, nil
}

func NewOrderItems(items []OrderItemInput, productPrices map[string]float64) ([]*OrderItem, error) {
	if len(items) < 1 || len(productPrices) < 1 {
		return nil, apperrors.ErrInvalidInput
	}

	orderItems := []*OrderItem{}
	for _, item := range items {
		price, ok := productPrices[item.ProductID]
		if !ok {
			return nil, apperrors.ErrInvalidInput
		}

		orderItems = append(orderItems, &OrderItem{
			MenuItemID: item.ProductID,
			Quantity:   item.Quantity,
			Price:      price,
		})
	}
	return orderItems, nil
}
