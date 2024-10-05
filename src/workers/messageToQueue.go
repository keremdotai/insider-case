package workers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"kerem.ai/insider/database"
)

func MessageToQueue() {
	// Set the sleep time for message generation
	valueStr := os.Getenv("SLEEP_TIME_FOR_MESSAGE_QUEUING")
	valueFloat, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		panic(fmt.Sprintf("cannot parse the sleep time for message queing : %v", valueFloat))
	}

	if valueFloat < 0.2 || valueFloat > 1 {
		panic("sleep time for message queuing should be between 0.2 and 1")
	}

	sleepTime := int(valueFloat * 1000)

	// Initialize queue with the messages that were in the queue before the system was stopped
	err = insertMessageToQueue(true)
	if err != nil {
		panic(fmt.Sprintf("cannot insert the messages that were in the queue before the system was stopped: %v", err))
	}

	// Generate messages
	for {
		// Generate a message
		err := insertMessageToQueue(false)
		if err != nil {
			continue
		}

		// Sleep
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}

func insertMessageToQueue(inQueue bool) error {
	// Find messages and insert them to the queue
	messages, err := database.FindMessages(inQueue)
	if err != nil {
		return err
	}

	for _, message := range messages {
		if len(messageQueue) == cap(messageQueue) {
			return nil
		}

		// Update the message as in the queue
		err = database.UpdateMessageInQueueStatus(message.ID)
		if err != nil {
			return err
		}

		// Insert the message to the queue
		messageQueue <- message
	}

	return nil
}
