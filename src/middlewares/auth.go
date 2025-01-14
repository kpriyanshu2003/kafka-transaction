package middleware

import (
	"example.com/kafka/src/config"
	"github.com/gofiber/fiber/v2"
)

func AuthToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token Required"})
		}

		db := config.DB.Db
		user, err := db.Exec("select * from users where email = $1", authHeader)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid Token"})
		}
		c.Locals("user", user)
		return c.Next()
	}
}
