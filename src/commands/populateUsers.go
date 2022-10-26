package main

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/bxcodec/faker/v3"
)

// docker-compose exec backend shで接続
// go run src/commands/populateUsers.go で実行
func main() {
	database.Connect()

	for i := 0; i < 30; i++ {
		seller := models.User{
			FirstName:    faker.FirstName(),
			LastName:     faker.LastName(),
			Email:        faker.Email(),
			IsSeller: true,
		}

		sellers.SetPassword("1234")
		database.DB.Create(&seller)
	}
}
