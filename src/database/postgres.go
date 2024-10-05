package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"kerem.ai/insider/models"
)

// Define PostgresDB as a pointer to gorm.DB
var PostgresDB *gorm.DB

// Connect to the Postgres database
func ConnectPostgresDB() {
	// Get environment variables for Postgres connection
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DBNAME")

	// Create a DSN string for Postgres connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port,
	)

	// Connect to the Postgres database
	var err error
	PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

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
func UpdateMessageSentStatus(messageID uint) error {
	return PostgresDB.Model(&models.Message{}).Where("id = ?", messageID).Update("status", true).Error
}

// Update the message as in the queue
func UpdateMessageInQueueStatus(messageID uint) error {
	return PostgresDB.Model(&models.Message{}).Where("id = ?", messageID).Update("in_queue", true).Error
}

// Message find operation
func FindMessages(inQueue bool) ([]models.Message, error) {
	var messages []models.Message
	var message models.Message
	var query *gorm.DB

	if inQueue {
		// Find all messages that are in the queue and not sent
		query = PostgresDB.Where("in_queue = ? AND sent = ?", true, false).Find(&messages)
	} else {
		// Find one message that is not in the queue and not sent
		query = PostgresDB.Where("in_queue = ? AND sent = ?", false, false).First(&message)
		messages = append(messages, message)
	}

	if query.Error != nil {
		return nil, query.Error
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
