package handlers

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"kerem.ai/insider/workers"
)

var isActive bool = true // Default true, message sending starts when the server starts
var mu sync.Mutex

// Handler for `/startAndStop` endpoint
func StartAndStopMessageSending(c echo.Context) error {
	// Lock the mutex
	mu.Lock()
	defer mu.Unlock()

	if isActive {
		workers.StopMessageSenderWorkers() // Stop the message sender workers
		isActive = false

		return c.JSON(
			http.StatusOK,
			map[string]string{"message": "Message sending stopped"},
		)
	}

	go workers.StartMessageSenderWorkers() // Start the message sender workers
	isActive = true

	return c.JSON(
		http.StatusOK,
		map[string]string{"message": "Message sending started"},
	)
}
