package controllers

import (
	"fmt"

	"example.com/kafka/src/config"
	"example.com/kafka/src/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateTransaction(c *fiber.Ctx) error {
	payload := struct {
		Sender_id   uuid.UUID `json:"sender_id"`
		Receiver_id uuid.UUID `json:"receiver_id"`
		Amount      string    `json:"amount"`
		Type        string    `json:"type"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}

	message := models.Transaction{
		SendID: payload.Sender_id,
		RecvID: payload.Receiver_id,
		Amt:    payload.Amount,
		Status: "PENDING",
		Type:   payload.Type,
	}

	fmt.Println(message)
	return c.JSON(fiber.Map{"message": "Transaction Created"})
}

func GetBalance(c *fiber.Ctx) error {
	db := config.DB.Db
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "id is required"})
	}
	user, err := db.Exec("SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}
