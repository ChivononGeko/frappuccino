package router

import (
	"database/sql"
	"frappuchino/internal/handler"
	"frappuchino/internal/repository"
	"frappuchino/internal/service"
	"net/http"
)

func SetupRoutes(db *sql.DB) (*http.ServeMux, error) {
	inventRepo := repository.NewInventoryRepository(db)
	inventService := service.NewInventoryService(inventRepo)
	inventHandler := handler.NewInventHandler(inventService)

	menuRepo := repository.NewMenuRepository(db)
	menuService := service.NewMenuService(menuRepo, inventRepo)
	menuHandler := handler.NewMenuHandler(menuService)

	customerRepo := repository.NewCustomerRepository(db)

	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo, menuRepo, inventRepo, customerRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	reportRepo := repository.NewReportsRepository(db)
	serviceReports := service.NewReportsService(reportRepo)
	handlerReports := handler.NewReportsHandler(serviceReports)

	mux := http.NewServeMux()

	addRoutes(mux, "/inventory", InventoryRouter(inventHandler))
	addRoutes(mux, "/menu", MenuRouter(menuHandler))
	addRoutes(mux, "/orders", OrderRouter(orderHandler))
	addRoutes(mux, "/reports", ReportRouter(handlerReports))

	return mux, nil
}

func addRoutes(mux *http.ServeMux, path string, router http.Handler) {
	mux.Handle(path, router)
	mux.Handle(path+"/", router)
}
