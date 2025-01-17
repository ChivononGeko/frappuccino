package service

import (
	"encoding/json"
	"fmt"
	"frappuchino/internal/models"
	"log/slog"
	"time"
)

type OrderRepository interface {
	AddOrderRepository(order models.Order, orderItems []*models.OrderItem) error
	GetOrderRepository(id int) (*models.Order, error)
	GetAllOrdersRepository() ([]*models.Order, error)
	UpdateOrderRepository(id int, order models.Order, orderItems []*models.OrderItem) error
	DeleteOrderRepository(id int) error
	CloseOrderRepository(id int) error
	NumberOfOrderedItemsRepository(startDate, endDate time.Time) (map[string]int, error)
	AddOrdersRepository(orders []*models.Order, orderItems [][]*models.OrderItem) error
}

type InventRepo interface {
	UpdateInventoryForSale(quantities map[string]float64) error
}

type MenuRepo interface {
	GetMenuItemsAndPrice(productIDs []string) (map[string]float64, error)
	CalculateIngredientsForOrder(menuQuantities map[string]int) (map[string]float64, error)
}

type CustomerRepo interface {
	IndentCustomerID(customerName string, instructions json.RawMessage) (int, error)
}

type OrderService struct {
	orderRepo    OrderRepository
	menuRepo     MenuRepo
	inventRepo   InventRepo
	customerRepo CustomerRepo
}

func NewOrderService(oR OrderRepository, mR MenuRepo, iR InventRepo, cR CustomerRepo) *OrderService {
	return &OrderService{
		orderRepo:    oR,
		menuRepo:     mR,
		inventRepo:   iR,
		customerRepo: cR,
	}
}

func (s *OrderService) CreateOrderService(orderRequest models.CreateOrderRequest) error {
	order, orderItems, err := s.createObject(orderRequest)
	if err != nil {
		slog.Error("Service error in Create Order: creating object", "input item", orderRequest, "error", err)
		return err
	}

	err = s.orderRepo.AddOrderRepository(*order, orderItems)
	if err != nil {
		slog.Error("Service error in Create Order: adding objects", "order", order, "order items", orderItems, "error", err)
		return err
	}

	return nil
}

func (s *OrderService) GetAllOrdersService() ([]*models.Order, error) {
	orders, err := s.orderRepo.GetAllOrdersRepository()
	if err != nil {
		slog.Error("Service error in Get Orders: retrieving all order", "error", err)
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) GetOrderService(id int) (*models.Order, error) {
	order, err := s.orderRepo.GetOrderRepository(id)
	if err != nil {
		slog.Error("Service error in Get Order: retrieving order", "id", id, "error", err)
		return nil, err
	}
	return order, nil
}

func (s *OrderService) UpdateOrderService(id int, orderRequest models.CreateOrderRequest) error {
	order, orderItems, err := s.createObject(orderRequest)
	if err != nil {
		slog.Error("Service error in Update Order: failed to create object", "input item", orderRequest, "error", err)
		return err
	}

	err = s.orderRepo.UpdateOrderRepository(id, *order, orderItems)
	if err != nil {
		slog.Error("Service error in Update Order: failed to update objects", "id", id, "order", order, "order items", orderItems, "error", err)
		return err
	}

	return nil
}

func (s *OrderService) DeleteOrderService(id int) error {
	err := s.orderRepo.DeleteOrderRepository(id)
	if err != nil {
		slog.Error("Service error in Delete Order: deleting order", "id", id, "error", err)
		return err
	}

	return nil
}

func (s *OrderService) CloseOrderService(id int) error {
	err := s.orderRepo.CloseOrderRepository(id)
	if err != nil {
		slog.Error("Service error in Close Order: close order", "id", id, "error", err)
		return err
	}

	return nil
}

func (s *OrderService) AddOrdersService(ordersRequests []models.CreateOrderRequest) error {
	var orders []*models.Order
	var orderItemsLists [][]*models.OrderItem
	for _, orderRequest := range ordersRequests {
		order, orderItemsList, err := s.createObject(orderRequest)
		if err != nil {
			slog.Error("Service error in Create Orders: creating objects", "error", err)
			return err
		}
		orders = append(orders, order)
		orderItemsLists = append(orderItemsLists, orderItemsList)
	}

	if err := s.orderRepo.AddOrdersRepository(orders, orderItemsLists); err != nil {
		slog.Error("Service error in Ba Create Orders: adding orders", "error", err)
		return err
	}

	return nil
}

func (s *OrderService) createObject(orderRequest models.CreateOrderRequest) (*models.Order, []*models.OrderItem, error) {
	productPrices, totalAmount, err := s.validateOrder(orderRequest)
	if err != nil {
		slog.Error("Service error in create objects: failed to validate order", "order", orderRequest, "error", err)
		return nil, nil, err
	}

	customerId, err := s.customerRepo.IndentCustomerID(orderRequest.CustomerName, orderRequest.Instructions)
	if err != nil {
		slog.Error("Service error in create objects: failed to ident customer id", "order", orderRequest, "error", err)
		return nil, nil, err
	}

	order, err := models.NewOrder(customerId, totalAmount, orderRequest)
	if err != nil {
		slog.Error("Service error in create objects: failed to create order", "order", orderRequest, "error", err)
		return nil, nil, err
	}

	orderItems, err := models.NewOrderItems(orderRequest.Items, productPrices)
	if err != nil {
		slog.Error("Service error in create objects: failed to create order items", "order", orderRequest, "error", err)
		return nil, nil, err
	}

	return order, orderItems, nil
}

func (s *OrderService) validateOrder(order models.CreateOrderRequest) (map[string]float64, float64, error) {
	productIDs := make([]string, len(order.Items))
	quantitiesInOrder := make(map[string]int)
	for i, item := range order.Items {
		productIDs[i] = item.ProductID
		quantitiesInOrder[item.ProductID] += item.Quantity
	}

	menuItems, err := s.menuRepo.GetMenuItemsAndPrice(productIDs)
	if err != nil {
		slog.Error("Service error in validate order: failed to retrieve menu and prices", "error", err)
		return nil, 0, err
	}

	var totalAmount float64
	for _, item := range order.Items {
		price, exists := menuItems[item.ProductID]
		if !exists {
			slog.Error("Service error in validate order: item not exist in menu", "item ID", item.ProductID)
			return nil, 0, fmt.Errorf("product with ID %s not found in menu", item.ProductID)
		}
		totalAmount += price * float64(item.Quantity)
	}

	ingredientsRequired, err := s.menuRepo.CalculateIngredientsForOrder(quantitiesInOrder)
	if err != nil {
		slog.Error("Service error in validate order: failed to calculate ingredients", "quantities", quantitiesInOrder, "error", err)
		return nil, 0, err
	}

	if err := s.inventRepo.UpdateInventoryForSale(ingredientsRequired); err != nil {
		slog.Error("Service error in validate order: failed to update inventory", "error", err)
		return nil, 0, err
	}

	return menuItems, totalAmount, nil
}

func (s *OrderService) NumberOfOrderedItemsService(startDateStr, endDateStr string) (map[string]int, error) {
	if startDateStr == "" {
		startDateStr = "01.01.1900"
	}
	if endDateStr == "" {
		endDateStr = "12.12.2100"
	}

	layout := "02.01.2006"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		slog.Error("Handler error from Number of Ordered Items: invalid date format", "start date", startDateStr)
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		slog.Error("Handler error from Number of Ordered Items: invalid date format", "end date", endDateStr)
	}

	order, err := s.orderRepo.NumberOfOrderedItemsRepository(startDate, endDate)
	if err != nil {
		slog.Error("Handler error from Number of Ordered Items: failed retrieving number ordered items", "error", err)
		return nil, err
	}
	return order, nil
}
