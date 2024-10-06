package database

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"kerem.ai/insider/models"
)

// Define PostgresDB as a pointer to gorm.DB
var PostgresDB *gorm.DB

// Connect to the Postgres database
func ConnectPostgresDB() {
	// Get environment variables for Postgres connection
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	// Create a DSN string for Postgres connection
	dsn := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=5432",
		dbname, user, password, host,
	)

	// Connect to the Postgres database
	var err error
	PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Postgres: %v", err))
	}

	// Migrate the data model to the database
	models.MigratePostgresDB(PostgresDB)
}

// Insert a new message item to the database
func InsertMessage(recipent string, content string) (uint, error) {
	// Create a new message item
	message := models.Message{
		Recipent: recipent,
		Content:  content,
		Sent:     false,
		InQueue:  false,
	}

	// Insert the message item to the database
	if result := PostgresDB.Create(&message); result.Error != nil {
		return 0, result.Error
	}

	return message.ID, nil
}

// Update the message as sent
func UpdateMessageSentStatus(messageID uint, sent bool) error {
	return PostgresDB.Model(&models.Message{}).Where("id = ?", messageID).Update("sent", sent).Error
}

// Update the message as in the queue
func UpdateMessageInQueueStatus(messageID uint, inQueue bool) error {
	return PostgresDB.Model(&models.Message{}).Where("id = ?", messageID).Update("in_queue", inQueue).Error
}

// Message find operation
func FindMessages(sent bool, inQueue bool, firstOnly bool) ([]models.Message, error) {
	var messages []models.Message
	var message models.Message

	if firstOnly {
		// Find the first message with the given conditions
		if err := PostgresDB.Where("sent = ? AND in_queue = ?", sent, inQueue).First(&message).Error; err != nil {
			// Check if the record is not found
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return messages, nil
			} else {
				return nil, err
			}
		} else {
			messages = append(messages, message)
		}
	} else {
		// Find all messages with the given conditions
		if err := PostgresDB.Where("sent = ? AND in_queue = ?", sent, inQueue).Find(&messages).Error; err != nil {
			return nil, err
		}
	}

	if inQueue {
		// Find all messages that are in the queue and not sent
		if err := PostgresDB.Where("in_queue = ? AND sent = ?", true, false).Find(&messages).Error; err != nil {
			return nil, err
		}
	} else {
		if err := PostgresDB.Where("in_queue = ? AND sent = ?", false, false).First(&message).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return messages, nil
			} else {
				return nil, err
			}
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// Close the connection to the Postgres database
func ClosePostgresDBConnection() error {
	// Get the underlying database connection
	db, err := PostgresDB.DB()
	if err != nil {
		return err
	}

	// Close the database connection
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
