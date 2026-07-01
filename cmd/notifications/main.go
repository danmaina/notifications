package main

import (
	"context"
	"github.com/danmaina/logger"
	"messaging/internal/api"
	"messaging/internal/configs"
	"github.com/danmaina/infra/v2/rabbitmq"
	"messaging/internal/mailer"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config, err := configs.ReadConfigs()
	if err != nil {
		logger.FATAL("Could Not Initialize Configs: ", err)
	}

	logger.SetLogLevel(config.ApplicationConfigs.LogLevel)

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.INFO("Received shutdown signal, initiating graceful shutdown...")
		cancel()
	}()

	// Start RabbitMQ Consumer in a goroutine
	go func() {
		logger.INFO("Starting RabbitMQ Consumer...")
		rabbitConsumer := &rabbitmq.RabbitMQConsumer{
			URL:       config.RabbitMQ.URL,
			QueueName: "NOTIFICATIONS",
			Handler: func(body []byte) error {
				return mailer.ProcessEmailRequest(body, config)
			},
		}
		err := rabbitConsumer.Start(ctx)
		if err != nil {
			logger.FATAL("RabbitMQ Consumer failed: ", err)
		}
	}()

	// Start HTTP Server for /health endpoint
	mux := http.NewServeMux()
	mux.HandleFunc("/health", api.HealthHandler)

	servePort := ":" + config.ApplicationConfigs.Port
	logger.INFO("Starting Health Check HTTP Server on port ", servePort)
	
	// Create an HTTP server and run it
	server := &http.Server{
		Addr:    servePort,
		Handler: mux,
	}

	go func() {
		errServe := server.ListenAndServe()
		if errServe != nil && errServe != http.ErrServerClosed {
			logger.FATAL("Could Not start health server on port ", servePort, " associated error: ", errServe)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	logger.INFO("Shutting down health server...")
	server.Shutdown(context.Background())
	logger.INFO("Application successfully stopped")
}
