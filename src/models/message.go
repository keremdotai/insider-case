package models

import (
	"gorm.io/gorm"
)

// Data model for the message item in database
type Message struct {
	ID       uint   `gorm:"primaryKey" json:"-"`
	Recipent string `validate:"len=14" json:"to"`    // Phone number must be in the form of +905xxxxxxxxx
	Content  string `gorm:"size:140" json:"content"` // Size limit for message content
	Sent     bool   `gorm:"default:false" json:"-"`  // Whether the message is sent or not
	InQueue  bool   `gorm:"default:false" json:"-"`  // Whether the message is in the queue or not
}

// Migrate data model to database
func MigratePostgresDB(db *gorm.DB) {
	db.AutoMigrate(&Message{})
}
