package main

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"

	"github.com/go-redis/redis/v8"
)

func main() {
	var users []models.User

	database.Connect()
	database.SetupRedis()

	ctx := context.Background()

	database.DB.Find(&users, models.User{
		IsSeller: true,
	})

	for _, user := range users{
		seller := models.Seller(user)
		seller.CalculateRevenue(database.DB)

		database.Cache.ZAdd(ctx, "rankings", &redis.Z{
			Score: *seller.Revenue,
			Member: user.Name(),
		})
	}
}