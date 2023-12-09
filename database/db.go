package database

import (
	"felix1234567890/go-trello/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectToDB() {
	user, password, database := os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DATABASE")
	fmt.Println(user, password, database)
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to database", err.Error())
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Cannot migrate", err.Error())
	}
	fmt.Println("Connected to database")

}
