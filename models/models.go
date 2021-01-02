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
	if tx := DB.Raw("DELETE FROM rooms;"); tx.Error != nil {
		return tx.Error
	}
	if tx := DB.Raw("DELETE FROM messages;"); tx.Error != nil {
		return tx.Error
	}
	if tx := DB.Create(&Room{Name: "Test Room"}); tx.Error != nil {
		return tx.Error
	}
	return nil
}
