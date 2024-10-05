package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"kerem.ai/insider/database"
	"kerem.ai/insider/models"
)

var webhookUrl string = os.Getenv("WEBHOOK_URL")  // Webhook URL
var messageQueue = make(chan models.Message, 100) // Message queue
var wg sync.WaitGroup                             // Wait group for workers
var stopSignal = make(chan bool)                  // Stop signal for workers

// Start the message sender workers
func StartMessageSenderWorkers() {
	wg.Add(2)

	go messageSenderWorker()
	go messageSenderWorker()

	wg.Wait() // Wait for the workers to finish
}

// Message sender worker
func messageSenderWorker() {
	defer wg.Done()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	client := &http.Client{}

	for {
		select {
		// Waiting 2 minutes for the next message
		case <-ticker.C:
			message, ok := <-messageQueue

			// Check if the message queue is closed
			if !ok {
				return
			}

			// Send the message to the webhook
			err := sendRequestToWebhook(client, &message)
			if err != nil {
				messageQueue <- message
			}
		// Receiving stop signal
		case <-stopSignal:
			return
		}
	}
}

// Add message to the queue
func AddMessageToQueue(message models.Message) {
	messageQueue <- message
}

// Close the message queue
func CloseMessageQueue() {
	close(messageQueue)
}

// Stop the message sender workers
func StopMessageSenderWorkers() {
	stopSignal <- true
	stopSignal <- true
}

// Send request to the webhook
func sendRequestToWebhook(client *http.Client, message *models.Message) error {
	// Prepare the payload for the webhook request
	payload := models.Payload{To: message.Recipent, Content: message.Content}
	data, _ := json.Marshal(payload)

	// Create a new request
	req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", "INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Read the response body
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to send message: %v", resp.Status)
	}

	// Update the message status
	err = database.UpdateMessageSentStatus(message.ID)
	if err != nil {
		return err
	}

	var responseData map[string]string // Response data
	json.NewDecoder(resp.Body).Decode(&responseData)

	messageID := responseData["messageId"]
	err = database.AppendMessageID("sent_messages", messageID)
	if err != nil {
		return err
	}

	return nil
}
