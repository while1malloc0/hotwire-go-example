package models

import "gorm.io/gorm"

// Message represents a chat message inside of a Room
type Message struct {
	gorm.Model

	Room   Room
	RoomID int

	Content string
}

// CreateMessage persists a new Message to the Database
func CreateMessage(message *Message) error {
	tx := DB.Create(message)
	return tx.Error
}
