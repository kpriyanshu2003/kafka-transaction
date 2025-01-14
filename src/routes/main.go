package routes

import "github.com/gofiber/fiber/v2"

func Routes(r fiber.Router) {
	r.Route("/transaction", TransactionRoutes)
}
