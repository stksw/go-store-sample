package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"ambassador/src/middlewares"
	"strconv"
	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
)

func Link(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var links []models.Link
	database.DB.Where("user_id = ?", id).Find(&links)

	for i, link := range links {
		var orders []models.Order
		database.DB.Where("code = ? and complete = true", link.Code).Find(&orders)
		links[i].Orders = orders
	}

	return c.JSON(links)
}

type CreateLinkRequest struct {
	Products []int
}

func CreateLink(c *fiber.Ctx) error {
	var request CreateLinkRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	id, _ := middlewares.GetUserId(c)
	link := models.Link {
		UserId: id,
		Code: faker.Username(),
	}

	for _, productId := range request.Products {
		product := models.Product{}
		product.Id = uint(productId)
		link.Products = append(link.Products, product)
	}

	database.DB.Create(&link)
	return c.JSON(link)
}

func Stats(c *fiber.Ctx) error {
	var links []models.Link
	var orders []models.Order
	var result []interface{}

	id, _ := middlewares.GetUserId(c)
	database.DB.Find(&links, models.Link{
		UserId: id,
	})

	for _, link := range links {
		database.DB.Preload("OrderItems").Find(&orders, &models.Order{
			Code: link.Code,
			Complete: true,
		})

		revenue := 0.0
		for _, order := range orders {
			revenue += order.GetTotal()
		}

		result = append(result, fiber.Map{
			"code": link.Code,
			"count": len(orders),
			"revenue": revenue,
		})
	}
	return c.JSON(result)
}

func GetLink(c *fiber.Ctx) error {
	var link models.Link
	code := c.Params("code")

	database.DB.Preload("User").Preload("Products").First(&link, &models.Link{
		Code: code,
	})
	return c.JSON(link)
}