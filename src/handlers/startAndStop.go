package handlers

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"kerem.ai/insider/models"
	"kerem.ai/insider/workers"
)

var isActive bool = true // Default true, message sending starts when the server starts
var mu sync.Mutex

// @Summary Start or stop the message sending.
// @Description It enables you to start or stop the message sending according to "action" parameter.
// @Accept  json
// @Produce  json
// @Param  request body models.StartAndStopRequest true "Action to start or stop"
// @Success 200 {object} map[string]string "Message describing the status of the message sending service"
// @Failure 400 {object} map[string]string "Error message explaining the invalid request payload or action"
// @Router /start-stop [post]
func StartAndStopMessageSending(c echo.Context) error {
	// Lock the mutex
	mu.Lock()
	defer mu.Unlock()

	// Bind the incoming JSON body to the struct and validate it
	var req models.StartAndStopRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	if req.Action != "start" && req.Action != "stop" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action. Action must be either 'start' or 'stop'"})
	}

	// Handle action
	if req.Action == "start" {
		// Check if the message sending is already active
		if isActive {
			return c.JSON(
				http.StatusOK,
				map[string]string{"message": "Message sending is already active"},
			)
		} else {
			// Start the message sender workers
			go workers.StartMessageSenderWorkers()
			isActive = true // Set the message sending status to true

			return c.JSON(
				http.StatusOK,
				map[string]string{"message": "Message sending started"},
			)
		}
	} else {
		// Check if the message sending is already stopped
		if isActive {
			// Stop the message sender workers
			workers.StopMessageSenderWorkers()
			isActive = false // Set the message sending status to false

			return c.JSON(
				http.StatusOK,
				map[string]string{"message": "Message sending stopped"},
			)
		} else {
			return c.JSON(
				http.StatusOK,
				map[string]string{"message": "Message sending is already stopped"},
			)
		}
	}
}
