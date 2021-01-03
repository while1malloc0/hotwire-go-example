package models

import (
	"gorm.io/gorm"
)

// Room represents a chat room
type Room struct {
	gorm.Model

	Name     string
	Messages []Message
}

// CreateRoom persists a new Room with the given name to the database
func CreateRoom(name string) error {
	tx := DB.Create(&Room{Name: name})
	return tx.Error
}

// FindRoom attempts to lookup a room by its ID
func FindRoom(id uint64) (*Room, error) {
	var room Room
	tx := DB.First(&room, id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	DB.Model(&room).Association("Messages").Find(&room.Messages)

	return &room, nil
}

// ListRooms retrieves a list of all Rooms from the database
func ListRooms() ([]*Room, error) {
	rows, err := DB.Model(&Room{}).Rows()
	if err != nil {
		return nil, err
	}

	var rooms []*Room
	for rows.Next() {
		var room Room
		DB.ScanRows(rows, &room)
		rooms = append(rooms, &room)
	}

	return rooms, nil
}

// UpdateRoom makes changes to an existing Room given a map of column:value pairs
// e.g. map[string]interface{}{"name": "some new name"}
func UpdateRoom(room *Room, updates map[string]interface{}) error {
	tx := DB.First(room, room.ID).Updates(updates)
	return tx.Error
}

// DeleteRoom deletes an existing Room
func DeleteRoom(room *Room) error {
	tx := DB.Delete(room)
	return tx.Error
}
