package db

import (
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Get() *gorm.DB {
	if db != nil {
		return db
	}

	db, err := gorm.Open(sqlite.Open("database.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database", err)
	}

	return db
}

func Directory() string {
	return "$XDG_DATA_HOME"
}
