package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"kerem.ai/insider/database"
)

// @Summary List all sent messages.
// @Description Retrieves a list of all messages that have been sent.
// @Produce json
// @Success 200 {array} models.Message "A list of sent messages"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /sent-messages [get]
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
