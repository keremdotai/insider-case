package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"kerem.ai/insider/database"
)

func ListSentMessages(c echo.Context) error {
	// Find all messages that are sent
	messages, err := database.FindMessages(true, true, false)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		messages,
	)
}
