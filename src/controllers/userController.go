package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func Sellers(c *fiber.Ctx) error {
	var users []models.User

	database.DB.Where("is_seller = true").Find(&users)
	return c.JSON(users)
}

func Rankings(c *fiber.Ctx) error {
	rankings, err := database.Cache.ZRevRangeByScoreWithScores(context.Background(), "rankings", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()
	if err != nil {
		return err
	}

	result := make(map[string]float64)
	for _, ranking := range rankings{
		fmt.Println(ranking.Member)
		result[ranking.Member.(string)] = ranking.Score
	}

	return c.JSON(result)
}
