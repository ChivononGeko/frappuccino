package models

import (
	"encoding/json"
	"frappuchino/internal/apperrors"
)

type CreateOrderRequest struct {
	CustomerName  string           `json:"customer_name"`
	PaymentMethod string           `json:"payment_method"`
	Items         []OrderItemInput `json:"items"`
	Instructions  json.RawMessage  `json:"instructions"`
}

type OrderItemInput struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

func NewCreateOrder(createOrder CreateOrderRequest) (*CreateOrderRequest, error) {
	if createOrder.CustomerName == "" || createOrder.PaymentMethod == "" {
		return nil, apperrors.ErrInvalidInput
	}

	for _, createcreateOrderItem := range createOrder.Items {
		if createcreateOrderItem.ProductID == "" || createcreateOrderItem.Quantity <= 0 {
			return nil, apperrors.ErrInvalidInput
		}
	}

	return &CreateOrderRequest{
		CustomerName:  createOrder.CustomerName,
		PaymentMethod: createOrder.PaymentMethod,
		Items:         createOrder.Items,
		Instructions:  createOrder.Instructions,
	}, nil
}
