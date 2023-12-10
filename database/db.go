package database

import (
	"felix1234567890/go-trello/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	user, password, database := os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to database", err.Error())
	}

	DB = db
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Cannot migrate", err.Error())
	}
	connection, err := db.DB()
	if err != nil {
		log.Fatal("Cannot get connection", err.Error())
	}
	connection.SetMaxIdleConns(5)
	connection.SetMaxOpenConns(10)

	fmt.Println("Connected to database")

}
