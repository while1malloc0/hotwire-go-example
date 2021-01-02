package models

import (
	"gorm.io/gorm"
)

type Room struct {
	gorm.Model

	Name     string
	Messages []Message
}

func CreateRoom(name string) error {
	tx := DB.Create(&Room{Name: name})
	return tx.Error
}

func FindRoom(id uint64) (*Room, error) {
	var room Room
	tx := DB.First(&room, id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	DB.Model(&room).Association("Messages").Find(&room.Messages)

	return &room, nil
}

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

func UpdateRoom(room *Room, updates map[string]interface{}) error {
	tx := DB.First(room, room.ID).Updates(updates)
	return tx.Error
}
