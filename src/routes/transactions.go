package routes

import (
	"example.com/kafka/src/controllers"
	"github.com/gofiber/fiber/v2"
)

func TransactionRoutes(r fiber.Router) {
	r.Post("/", controllers.CreateTransaction)
	r.Get("/:id?", controllers.GetBalance)
}
