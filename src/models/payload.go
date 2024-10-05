package models

// Data model for the payload for webhook requests
type Payload struct {
	To      string `json:"to"`      // Phone number
	Content string `json:"content"` // Message content
}
