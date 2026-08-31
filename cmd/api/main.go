package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SaSiaing/backend/internal/database"
	"github.com/SaSiaing/backend/internal/health"
	"github.com/labstack/echo/v5"
)

const (
	defaultPort     = "8080"
	databaseTimeout = 5 * time.Second
	gracefulTimeout = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	databaseContext, cancel := context.WithTimeout(ctx, databaseTimeout)
	defer cancel()

	pool, err := database.Open(databaseContext, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer pool.Close()

	e := echo.New()
	e.GET("/health", health.Handler(pool))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	slog.Info("database connected")
	server := echo.StartConfig{
		Address:         ":" + port,
		GracefulTimeout: gracefulTimeout,
	}
	if err := server.Start(ctx, e); err != nil {
		return fmt.Errorf("serve HTTP: %w", err)
	}

	return nil
}
