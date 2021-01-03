package models

import (
	"errors"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is a connection to the Database
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

// Migrate runs migrations for all models.
func Migrate() error {
	return DB.AutoMigrate(Room{}, Message{})
}

// IsRecordNotFound determines if an error was caused by a database transaction
// returning no records
func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// Seed creates test records in the database
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
