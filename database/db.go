package database

import (
	"felix1234567890/go-trello/models"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectToDB initializes and returns a connection to the MySQL database using GORM.
// It reads database configuration (user, password, database name) from environment variables.
// It also performs an auto-migration for the User model.
// Returns a pointer to the gorm.DB instance and an error if any step fails.
func ConnectToDB() (*gorm.DB, error) {
	user, password, database := os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, fmt.Errorf("cannot migrate: %w", err)
	}
	connection, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("cannot get connection: %w", err)
	}
	connection.SetMaxIdleConns(5)
	connection.SetMaxOpenConns(10)

	fmt.Println("Connected to database")
	return db, nil
}
