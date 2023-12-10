package utils

import (
	"errors"
	"felix1234567890/go-trello/database"
	"felix1234567890/go-trello/models"
	"fmt"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/go-playground/validator"
)

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

func ValidateRequest(data interface{}) error {
	validate := validator.New()
	err := validate.Struct(data)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return errors.New(fmt.Sprintf("Field: %s, Error: %s\n", e.Field(), e.Tag()))
		}
	}
	return nil
}
