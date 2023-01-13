package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"pArtour/go-blog-api/pkg/config"
	"syscall"
	"time"
)

func main() {
	// Create logger.
	prLogger, _ := zap.NewProduction()
	defer prLogger.Sync()
	logger := prLogger.Sugar()

	logger.Info("Hello World")

	cfg, err := config.NewConfig()

	if err != nil {
		logger.Fatalw("Failed to parse and initialize config", "err", err)
	}

	// Print configuration.
	fmt.Printf("Configuration: %+v\n", cfg)

	// Create http server to use the router.
	server := &http.Server{
		Addr:         cfg.HttpAddress,
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		logger.Infow("Starting server", "address", cfg.HttpAddress)
		if err := server.ListenAndServe(); err != nil {
			logger.Fatalw("Failed to start server", "err", err, "http_address", cfg.HttpAddress)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	<-signals

	ctx, cancel := context.WithTimeout(context.TODO(), 15*time.Second)

	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		logger.Fatalw("Failed to shutdown server", "err", err, "http_address", cfg.HttpAddress)
	}

	logger.Infow("Server gracefully stopped", "http_address", cfg.HttpAddress)
}
