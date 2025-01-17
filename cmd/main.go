package main

import (
	"frappuchino/internal/config"
	"frappuchino/internal/db"
	"frappuchino/internal/router"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Configuration loading error", "error", err)
		os.Exit(1)
	}
	slog.Info("Configuration loaded successfully")

	dataBase, err := db.InitDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		slog.Error("Database connection failed", "host", cfg.DBHost, "port", cfg.DBPort, "error", err)
		os.Exit(1)
	}
	defer dataBase.Close()
	slog.Info("Database connection successfully")

	mux, err := router.SetupRoutes(dataBase)
	if err != nil {
		slog.Error("Failed to set up routes", "error", err)
		os.Exit(1)
	}
	slog.Info("Setap router successfully")

	slog.Info("Starting server", "port", cfg.APIPort)
	if err := http.ListenAndServe(":"+cfg.APIPort, mux); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
