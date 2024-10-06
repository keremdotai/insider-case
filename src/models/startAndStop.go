package models

// Data model for `/start-stop` endpoint
type StartAndStopRequest struct {
	Action string `json:"action"` // Action to be taken, either "start" or "stop"
}
