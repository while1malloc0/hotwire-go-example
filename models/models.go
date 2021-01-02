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
