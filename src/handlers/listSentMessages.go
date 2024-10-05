package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"kerem.ai/insider/database"
)

func ListSentMessages(c echo.Context) error {
	sentMessageIDs, err := database.RetrieveListItems("sent_messages")
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]interface{}{"sent_messages": sentMessageIDs},
	)
}
