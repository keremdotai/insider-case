package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"kerem.ai/insider/database"
	"kerem.ai/insider/handlers"
	"kerem.ai/insider/workers"
)

func main() {
	// Get host and port from the environment variables
	port := os.Getenv("PORT")

	if port == "" {
		panic("port cannot be empty")
	}

	serverAddress := fmt.Sprintf("0.0.0.0:%s", port)

	// Initialize the connection to the databases
	// If any of the connections fail, the program will panic
	database.ConnectPostgresDB()
	database.ConnectRedisDB()

	// Start message generation and queuing
	go workers.GenerateMessages()
	go workers.MessageToQueue()

	// Start the message sender workers
	go workers.StartMessageSenderWorkers()

	// Define server and its routes
	e := echo.New()
	e.Logger.SetLevel(log.INFO)

	e.GET("/sent-messages", handlers.ListSentMessages)
	e.POST("/start-stop", handlers.StartAndStopMessageSending)

	// Start the server
	go func() {
		err := e.Start(serverAddress)
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Capture shutdown signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait till the signal is received
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop workers and close the database connections
	workers.StopMessageSenderWorkers()
	workers.CloseMessageQueue()

	if err := database.ClosePostgresDBConnection(); err != nil {
		log.Fatalf("Error shutting down Postgres: %v", err)
	}

	if err := database.CloseRedisDBConnection(); err != nil {
		log.Fatalf("Error shutting down Redis: %v", err)
	}

	// Shutdown the server
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Info("Server shutdown completed")
}
