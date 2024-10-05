package workers

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/icrowley/fake"
	"kerem.ai/insider/database"
)

func GenerateMessages() {
	// Set the sleep time for message generation
	valueStr := os.Getenv("SLEEP_TIME_FOR_MESSAGE_GENERATION")
	valueFloat, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		panic(fmt.Sprintf("cannot parse the sleep time for message generation: %v", valueFloat))
	}

	if valueFloat < 0.2 || valueFloat > 1 {
		panic("sleep time for message generation should be between 0.2 and 1")
	}

	sleepTime := int(valueFloat * 1000)

	// Generate messages
	for {
		// Generate a message
		err := generateMessage()
		if err != nil {
			continue
		}

		// Sleep
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)
	}
}

func generateMessage() error {
	// Generate a random phone number
	recipent := fmt.Sprintf("+905%v", fake.DigitsN(9))

	// Generate a random message
	content := fake.Sentence()
	if len(content) > 140 {
		content = content[:140]
	}

	// Insert the message to the database
	_, err := database.InsertMessage(recipent, content)
	if err != nil {
		return fmt.Errorf("cannot insert the generated message: %v", err)
	}

	return nil
}
