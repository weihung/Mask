package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	SqlDb *gorm.DB
)

func init() {
	db, err := gorm.Open("mysql", sqlServer)
	if err != nil {
		log.Fatal(err)
	}
	SqlDb = db
	db.AutoMigrate(&Left{})
	db.AutoMigrate(&Store{})
	log.Println("SQL server auto migrate completed")
}

// CloseDB : Close DB
func CloseDB() {
	SqlDb.Close()

	log.Println("Close database finish")
}
