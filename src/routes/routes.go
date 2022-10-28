package routes

import (
	"ambassador/src/controllers"
	"ambassador/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	api.Get("products", controllers.ProductsFrontend)
	api.Get("search_products", controllers.SearchProducts)

	
	// seller api
	seller := api.Group("seller")
	seller.Post("register", controllers.Register)
	seller.Post("login", controllers.Login)

	sellerAuthenticated := seller.Use(middlewares.IsAuthenticated)
	sellerAuthenticated.Post("logout", controllers.Logout)
	sellerAuthenticated.Get("profile", controllers.Profile)
	sellerAuthenticated.Put("users/info", controllers.UpdateInfo)
	sellerAuthenticated.Put("users/password", controllers.UpdatePassword)

	// admin api
	admin := api.Group("admin")
	admin.Post("register", controllers.Register)
	admin.Post("login", controllers.Login)

	adminAuthenticated := admin.Use(middlewares.IsAuthenticated)
	adminAuthenticated.Post("logout", controllers.Logout)
	adminAuthenticated.Get("profile", controllers.Profile)
	adminAuthenticated.Put("users/info", controllers.UpdateInfo)
	adminAuthenticated.Put("users/password", controllers.UpdatePassword)
	adminAuthenticated.Get("sellers", controllers.Sellers)
	adminAuthenticated.Get("products", controllers.Products)
	adminAuthenticated.Get("products/:id", controllers.GetProduct)
	adminAuthenticated.Post("products", controllers.CreateProduct)
	adminAuthenticated.Delete("products/:id", controllers.DeleteProduct)
	adminAuthenticated.Get("users/:id/links", controllers.Link)
	adminAuthenticated.Get("orders", controllers.Orders)
}
