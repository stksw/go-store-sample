package main

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"math/rand"

	"github.com/bxcodec/faker/v3"
)

// docker-compose exec backend shで接続
// go run src/commands/populateUsers.go で実行
func main() {
	database.Connect()

	for i := 0; i < 30; i++ {
		var orderItems []models.OrderItem

		for j := 0; j < rand.Intn(5); j++ {
			price := float64(rand.Intn(90) + 10)
			qty := uint(rand.Intn(5))

			orderItems = append(orderItems, models.OrderItem{
				ProductTitle:      faker.Word(),
				Price:             price,
				Quantity:          qty,
				AdminRevenue:      0.9 * price * float64(qty),
				AmbassadorRevenue: 0.1 * price * float64(qty),
			})
		}
		order := models.Order{
			UserId:          uint(rand.Intn(30) + 1),
			Code:            faker.Username(),
			AmbassadorEmail: faker.Email(),
			FirstName:       faker.FirstName(),
			LastName:        faker.LastName(),
			Email:           faker.Email(),
			Complete:        true,
			OrderItems:      orderItems,
		}

		database.DB.Create(&order)
	}
}
