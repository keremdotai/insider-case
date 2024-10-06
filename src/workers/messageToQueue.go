package workers

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"kerem.ai/insider/database"
)

func MessageToQueue() {
	// Set the sleep time for the error case
	sleepTime := 2000 // 2 seconds

	// Initialize queue with the messages that were in the queue before the system was stopped
	err := insertMessageToQueue(true, false)
	if err != nil {
		panic(fmt.Sprintf("cannot insert the messages that were in the queue before the system was stopped: %v", err))
	}

	// Generate messages
	for {
		// Generate a message
		err := insertMessageToQueue(false, true)
		if err != nil {
			log.Errorf("cannot insert the message to the queue: %v", err)
			time.Sleep(time.Millisecond * time.Duration(sleepTime))
		}
	}
}

func insertMessageToQueue(inQueue bool, firstOnly bool) error {
	// Check if the queue is full
	if len(messageQueue) == cap(messageQueue) {
		return nil
	}

	// Find applicable messages
	messages, err := database.FindMessages(false, inQueue, firstOnly)
	if err != nil {
		return err
	}

	for _, message := range messages {
		// Update the message as in the queue
		err = database.UpdateMessageInQueueStatus(message.ID, true)
		if err != nil {
			return err
		}

		// Insert the message to the queue
		messageQueue <- message
	}

	return nil
}
