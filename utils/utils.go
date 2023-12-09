package utils

import (
	"felix1234567890/go-trello/database"
	"felix1234567890/go-trello/models"
	"math/rand"

	"github.com/go-faker/faker/v4"
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
