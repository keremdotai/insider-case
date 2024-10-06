package workers

import (
	"fmt"
	"time"

	"github.com/icrowley/fake"
	"github.com/labstack/gommon/log"
	"kerem.ai/insider/database"
)

func GenerateMessages() {
	// Set the sleep time for message generation
	sleepTime := 200 // 0.2 seconds

	// Generate messages
	for {
		// Generate a message
		err := generateMessage()
		if err != nil {
			log.Errorf("cannot generate the message: %v", err)
			continue
		}

		// Sleep
		time.Sleep(time.Millisecond * time.Duration(sleepTime))
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
