package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func Sellers(c *fiber.Ctx) error {
	var users []models.User

	database.DB.Where("is_seller = true").Find(&users)
	return c.JSON(users)
}

// func Rankings(c *fiber.Ctx) error {
// 	rankings, err := database.Cache.ZRevRangeByScoreWithScores(context.Background(), "rankings", &redis.ZRangeBy{
// 		Min: "-inf",
// 		Max: "+inf",
// 	}).Result()

// 	result := make(map[string]float64)

// 	return c.JSON(result)
// }
