package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model

	Room   Room
	RoomID int

	Content string
}

func CreateMessage(message *Message) error {
	tx := DB.Create(message)
	return tx.Error
}
