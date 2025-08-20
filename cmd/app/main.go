// cmd/main.go
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crypto-observer/internal/api"
	"crypto-observer/internal/db"
	"crypto-observer/internal/service"
	"crypto-observer/pkg/config"
	"crypto-observer/pkg/logger"
)

func main() {
	// 1) конфиг + логгер
	cfg := config.MustLoad() // паникует, если конфига нет/битый
	logger.Init()            // читает уровень из cfg.Log.Level внутри
	log := logger.L()

	// 2) контекст отмены по сигналам
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// 3) storage (pgxpool) + авто-схема
	st, err := db.NewStorage(cfg.DB.DSN)
	if err != nil {
		log.WithError(err).Fatal("storage init failed")
	}
	defer st.Close()

	// 4) сервис
	svc := service.NewService(
		st,
		cfg.Collector.DefaultPeriodSeconds,
		cfg.Coingecko.BaseURL,
		time.Duration(cfg.Coingecko.TimeoutSec)*time.Second,
	)

	// 5) http router
	r := api.NewRouter(api.NewHandler(svc))

	// 6) http server
	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.WithField("addr", cfg.Server.Addr).Info("HTTP server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("HTTP server stopped")
		}
	}()

	// 7) блокирующее ожидание сигнала
	<-ctx.Done()
	log.Info("shutdown started")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)

	log.Info("shutdown complete")
	os.Exit(0)
}
