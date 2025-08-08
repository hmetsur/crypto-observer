// @title Crypto Observer API
// @version 1.0
// @description Микросервис для мониторинга цен криптовалют
// @host localhost:8080
// @BasePath /
package main

import (
	"net/http"
	"os"

	"crypto-observer/internal/api"
	"crypto-observer/internal/collector"
	"crypto-observer/internal/db"
	"crypto-observer/internal/logger"
	"crypto-observer/internal/service"

	_ "github.com/lib/pq"
)

func main() {
	logger.Init()

	dsn := os.Getenv("DB_DSN")
	logger.Log.WithField("dsn", dsn).Info("Starting app, connecting DB")

	storage, err := db.NewStorage(dsn)
	if err != nil {
		logger.Log.WithError(err).Fatal("DB connect failed")
	}
	logger.Log.Info("DB connected")

	coll := collector.NewCollector(storage, nil)
	srv := service.NewService(storage, coll)

	router := api.NewRouter(srv)

	addr := ":8080"
	logger.Log.WithField("addr", addr).Info("HTTP server listening")

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Log.WithError(err).Fatal("HTTP server stopped")
	}
}
