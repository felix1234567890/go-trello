// Package utils provides helper functions for various tasks within the application,
// such as request validation, password hashing, JWT creation, and standardized JSON responses.
package utils

import (
	"errors"
	"felix1234567890/go-trello/models"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var SECRET_KEY []byte

// Custom error types for better error handling
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
)

// InitSecretKey loads SECRET_KEY from environment after .env is loaded.
func InitSecretKey() {
	if os.Getenv("SECRET_KEY") == "" {
		log.Fatalf("FATAL: SECRET_KEY environment variable is not set.")
	}
	SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))
}

// FakeUserFactory creates a random number of fake users (between 5 and 10)
// and persists them to the database using the provided *gorm.DB connection.
// Usernames, emails, and passwords (unhashed) are generated using a faker library.
// Returns an error if any database operation fails.
func FakeUserFactory(db *gorm.DB) error {
	min := 5
	max := 10
	randomValue := rand.Intn(max-min) + min

	for i := 0; i < randomValue; i++ {
		plainPassword := faker.Password()
		hashedPassword, err := HashPassword(plainPassword)
		if err != nil {
			return fmt.Errorf("failed to hash password for fake user %d: %w", i, err)
		}
		user := models.User{
			Username: faker.Username(),
			Email:    faker.Email(),
			Password: hashedPassword, // Password is now properly hashed
		}
		result := db.Create(&user)
		if result.Error != nil {
			return fmt.Errorf("failed to create fake user %d: %w", i, result.Error)
		}
	}
	return nil
}

// ValidateRequest performs struct validation using "github.com/go-playground/validator".
// It takes any interface{} as input, which should be a pointer to a struct with validation tags.
// If validation errors occur, it returns a map[string]string where keys are field names
// and values are error messages.
// Returns nil if validation passes.
func ValidateRequest(data interface{}) map[string]string {
	if data == nil {
		return map[string]string{"request": "Request data cannot be nil"}
	}
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		validationErrors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors[e.Field()] = fmt.Sprintf("Field %s failed on the '%s' tag", e.Field(), e.Tag())
		}
		return validationErrors
	}
	return nil
}

// HandleErrorResponse sends a standardized JSON error response with a given status code and message.
// The JSON response format is `{"message": "error message"}`.
func HandleErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"message": message})
}

// JsonResponse sends a standardized JSON success response with a given status code and data.
func JsonResponse(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(data)
}

// HashPassword generates a bcrypt hash for the given password string.
// Returns the hashed password string and an error if hashing fails.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

// CheckPasswordHash compares a plain text password with a bcrypt hashed password.
// Returns nil if the password matches the hash, otherwise returns an error.
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

// CreateToken generates a new JWT token for a given user ID.
// The token is signed using HS256 and includes the user ID and an expiration time of 1 hour.
// SECRET_KEY (a package-level variable initialized from env) is used for signing.
// Returns the signed token string and an error if token generation or signing fails.
func CreateToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
	})
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
