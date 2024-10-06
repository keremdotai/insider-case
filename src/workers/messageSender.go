package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
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

	log.Info("All message sender workers are stopped")
}

// Message sender worker
func messageSenderWorker() {
	defer wg.Done()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	client := &http.Client{}

	for {
		select {
		// Waiting 2 seconds for the next message
		case <-ticker.C:
			log.Info("Waiting for the next message") //

			select {
			// Get the message from the queue
			case message, ok := <-messageQueue:
				// Check if the message queue is closed
				if !ok {
					log.Info("message queue is closed")
					return
				}

				log.Info("Message received to send")
				err := sendRequestToWebhook(client, &message)
				if err != nil {
					log.Errorf("failed to send message: %v", err)

					// Update the message in_queue status to false
					if err := database.UpdateMessageInQueueStatus(message.ID, false); err != nil {
						log.Errorf("failed to update message in_queue status: %v", err)
					}
				}

			default:
				log.Info("message queue is empty")
			}
		// Receiving stop signal
		case <-stopSignal:
			log.Info("Stopping message sender worker")
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
	dt := time.Now() // Timestamp for the request

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
	err = database.UpdateMessageSentStatus(message.ID, true)
	if err != nil {
		return err
	}

	// Parse response body
	var responseData map[string]string // Response data
	json.NewDecoder(resp.Body).Decode(&responseData)
	responseData["timestamp"] = dt.Format("2006-01-02 15:04:05")

	jsonItem, err := json.Marshal(responseData)
	if err != nil {
		return err
	}

	err = database.AppendSentMessage("sent_messages", string(jsonItem))
	if err != nil {
		return err
	}

	return nil
}
