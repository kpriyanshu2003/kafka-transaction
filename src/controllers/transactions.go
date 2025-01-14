package controllers

import (
	"encoding/json"
	"fmt"

	"example.com/kafka/src/config"
	"example.com/kafka/src/kafka/producer"
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
	// TODO: Validations

	message := models.Transaction{
		TransactionID: uuid.New(),
		SendID:        payload.Sender_id,
		RecvID:        payload.Receiver_id,
		Amt:           payload.Amount,
		Status:        "PENDING",
		Type:          payload.Type,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
	}

	err = producer.PushTransactionToQueue("wallet.pending", jsonMessage)
	if err != nil {
		fmt.Printf("Failed to push transaction to Kafka: %v\n", err)
	}

	return c.JSON(message)
}

func GetBalance(c *fiber.Ctx) error {
	db := config.DB.Db
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{"error": "id is required"})
	}
	var user models.User
	query := "SELECT * FROM users WHERE user_id = $1"
	err := db.QueryRow(query, id).Scan(&user.UserID, &user.Email, &user.Balance)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(&user)
}
