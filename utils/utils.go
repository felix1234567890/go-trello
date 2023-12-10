package utils

import (
	"errors"
	"felix1234567890/go-trello/database"
	"felix1234567890/go-trello/models"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = []byte(os.Getenv("SECRET_KEY"))

func FakeUserFactory() {
	min := 5
	max := 10
	randomValue := rand.Intn(max-min) + min

	for i := 0; i < randomValue; i++ {
		user := models.User{
			Username: faker.Username(),
			Email:    faker.Email(),
			Password: faker.Password(),
		}
		database.DB.Create(&user)
	}
}

func ValidateRequest(data interface{}) []error {
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		var validationErrors []error
		for _, e := range err.(validator.ValidationErrors) {
			errMsg := fmt.Sprintf("'%s' has a value of '%v' which does not satisfy '%s' constraint", e.Field(), e.Value(), e.Tag())
			validationErrors = append(validationErrors, errors.New(errMsg))
		}
		return validationErrors
	}
	return nil
}

func HandleErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"message": message})
}

func JsonResponse(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(data)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func CreateToken(id uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("Invalid Token")
	}

	return claims["id"].(string), nil
}
