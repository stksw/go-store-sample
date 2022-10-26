package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"

	"github.com/gofiber/fiber/v2"
)

func Orders(c *fiber.Ctx) error {
	var orders []models.Order

	database.DB.Preload("OrderItems").Find(&orders)
	for i, order := range orders {
		orders[i].Name = order.FullName()
		orders[i].Total = order.GetTotal()
	}

	return c.JSON(orders)
}

type CreateOrderRequest struct {
	Code      string
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
	City      string
	Zip       string
	Products  []map[string]int
}

func CreateOrder(c *fiber.Ctx) error {
	var request CreateOrderRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	link := models.Link{
		Code: request.Code,
	}
	database.DB.Preload("User").First(&link)

	if link.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Request",
		})
	}

	// order := models.Order{
	// 	Code:            link.Code,
	// 	UserId:          link.UserId,
	// 	AmbassadorEmail: link.User.Email,
	// 	FirstName:       request.FirstName,
	// 	LastName:        request.LastName,
	// 	Email:           request.Email,
	// 	Address:         request.Address,
	// 	Country:         request.Country,
	// 	City:            request.City,
	// 	Zip:             request.Zip,
	// }

	return c.JSON(link)

}
