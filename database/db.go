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
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "db"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Group{})
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
