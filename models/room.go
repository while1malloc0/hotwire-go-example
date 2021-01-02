package models

import (
	"strconv"

	"gorm.io/gorm"
)

type Room struct {
	gorm.Model

	Name     string
	Messages []Message
}

func FindRoom(idStr string) (*Room, error) {
	var room Room
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	room.ID = uint(id)
	tx := DB.Find(&room)
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

func UpdateRoom(id string, updates map[string]interface{}) error {
	tx := DB.First(&Room{}, id).Updates(updates)
	return tx.Error
}
