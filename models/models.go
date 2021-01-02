package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	if DB == nil {
		DB, err = gorm.Open(sqlite.Open("chat.db"), &gorm.Config{})
	}
	if err != nil {
		panic(err)
	}
}

func Migrate() error {
	return DB.AutoMigrate(Room{}, Message{})
}

func Seed() error {
	var room Room
	testRoomName := "Test Room"
	if DB.Where(&Room{Name: testRoomName}).First(&room); room.ID == 0 {
		if tx := DB.Create(&Room{Name: testRoomName}); tx.Error != nil {
			return tx.Error
		}
	}
	return nil
}
