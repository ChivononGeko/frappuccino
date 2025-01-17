package service

import (
	"fmt"
	"frappuchino/internal/models"
	"log/slog"
	"strconv"
)

type ReportsRepository interface {
	GetTotalSales() (*models.TotalPrice, error)
	GetPopularItems() ([]*models.PopularItem, error)
	SearchMenuItems(q string, minPrice, maxPrice float64) ([]map[string]interface{}, error)
	SearchOrders(q string, minPrice, maxPrice float64) ([]map[string]interface{}, error)
	OrderedItemByDayRepository(month string) (map[string]interface{}, error)
	OrderedItemByMonthRepository(year int) (map[string]interface{}, error)
}

type ReportsService struct {
	reportRepo ReportsRepository
}

func NewReportsService(or ReportsRepository) *ReportsService {
	return &ReportsService{
		reportRepo: or,
	}
}

func (s *ReportsService) TotalSalesReportService() (*models.TotalPrice, error) {
	totalSales, err := s.reportRepo.GetTotalSales()
	if err != nil {
		slog.Error("Service error in Total Sales: failed to get total sales", "error", err)
		return nil, err
	}

	return totalSales, nil
}

func (s *ReportsService) PopularItemsReportService() ([]*models.PopularItem, error) {
	popularItem, err := s.reportRepo.GetPopularItems()
	if err != nil {
		slog.Error("Service error in Total Sales: failed to get popular items", "error", err)
		return nil, err
	}
	return popularItem, nil
}

func (s *ReportsService) SearchService(q, filter, minPriceStr, maxPriceStr string) (map[string]interface{}, error) {
	var menuItems []map[string]interface{}
	var orders []map[string]interface{}
	var totalMatches int

	minPrice, err := strconv.ParseFloat(minPriceStr, 64)
	if err != nil {
		slog.Error("Service error from Search Service: failed parse minPrice to float", "minPrice", minPrice, "error", err)
		return nil, err
	}

	maxPrice, err := strconv.ParseFloat(maxPriceStr, 64)
	if err != nil {
		slog.Error("Service error from Search Service: failed parse maxPrice to float", "maxPrice", maxPrice, "error", err)
		return nil, err
	}

	if filter == "menu" || filter == "all" {
		menuItems, err = s.reportRepo.SearchMenuItems(q, minPrice, maxPrice)
		if err != nil {
			slog.Error("Service error from Search: failed retrieved menu items", "error", err)
			return nil, err
		}
		totalMatches += len(menuItems)
	}

	if filter == "orders" || filter == "all" {
		orders, err = s.reportRepo.SearchOrders(q, minPrice, maxPrice)
		if err != nil {
			slog.Error("Service error from Search: failed retrieved order items", "error", err)
			return nil, err
		}
		totalMatches += len(orders)
	}

	slog.Info("Search successfully", "total matches", totalMatches)
	return map[string]interface{}{
		"menu_items":    menuItems,
		"orders":        orders,
		"total_matches": totalMatches,
	}, nil
}

func (s *ReportsService) OrderedItemsByPeriodService(period, month, yearStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	var err error

	slog.Info("Values params", "period", period, "month", month, "year", yearStr)
	if period == "day" {
		if month == "" {
			slog.Error("Service error in Ordered Items by Period: missing month")
			return nil, fmt.Errorf("missing month")
		}

		if !checkMonth(month) {
			slog.Error("Service error in Ordered Items by Period: invalid month")
			return nil, fmt.Errorf("invalid month")
		}

		result, err = s.reportRepo.OrderedItemByDayRepository(month)
		if err != nil {
			slog.Error("Service error from Ordered Items by Period: failed retrieved ordered items by day", "month", month, "error", err)
			return nil, err
		}
	}

	if period == "month" {
		if yearStr == "" {
			slog.Error("Service error in Ordered Items by Period: missing year")
			return nil, fmt.Errorf("missing year")
		}
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			slog.Error("Service error in Ordered Items by Period: failed to parse year", "year", year, "error", err)
			return nil, err
		}
		result, err = s.reportRepo.OrderedItemByMonthRepository(year)
		if err != nil {
			slog.Error("Service error from Ordered Items by Period: failed retrieved ordered items by year", "year", month, "error", err)
			return nil, err
		}
	}

	return result, nil
}

func checkMonth(month string) bool {
	validMonth := map[string]bool{
		"january":   true,
		"february":  true,
		"march":     true,
		"april":     true,
		"may":       true,
		"june":      true,
		"july":      true,
		"august":    true,
		"september": true,
		"october":   true,
		"november":  true,
		"december":  true,
	}

	return validMonth[month]
}
