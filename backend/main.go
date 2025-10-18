package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dnswd/arus/config"
	"github.com/dnswd/arus/db"
	"github.com/dnswd/arus/health"
	"github.com/dnswd/arus/infra"
	"github.com/dnswd/arus/logger"
	"github.com/dnswd/arus/server"
	"github.com/dnswd/arus/user"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	// log.SetOutput(io.Discard)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Fprintln(os.Stdout, "app is starting")

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		return 1
	}

	if err := logger.InitWithFile(cfg.LogPath); err != nil {
		fmt.Fprintf(os.Stderr, "error initializing logger: %v\n", err)
		return 1
	}
	defer logger.Close()

	infra, err := infra.New(ctx, cfg)
	if err != nil {
		slog.InfoContext(ctx, "Failed to initialize infrastructure: %v", err)
		return 1
	}
	defer infra.Close()

	queries := db.New(infra.DB())

	// Wire up domain services
	userRepo := user.NewPostgresRepository(queries)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Health check
	healthHandler := health.NewHandler(infra.DB())

	// HTTP Server
	srv := server.New(healthHandler, userHandler /* , orderHandler */)
	if err := srv.Start(ctx, ":8080"); err != nil {
		slog.ErrorContext(ctx, "Failed to start server: %v", err)
		return 1
	}

	// Background scheduler (optional)
	// sched := scheduler.New(userService /* , orderService */)
	// if err := sched.Start(ctx); err != nil {
	//     log.Printf("Failed to start scheduler: %v", err)
	//     return 1
	// }

	// // Queue workers (optional)
	// work := worker.New(userService /* , orderService */)
	// if err := work.Start(ctx); err != nil {
	//     log.Printf("Failed to start workers: %v", err)
	//     return 1
	// }

	// Wait for interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	// Graceful shutdown
	slog.InfoContext(ctx, "Shutting down gracefully...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Stop(shutdownCtx); err != nil {
		slog.ErrorContext(ctx, "Server shutdown error: %v", err)
	}

	return 0
}
