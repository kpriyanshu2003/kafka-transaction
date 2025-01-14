package controllers

import (
	"example.com/kafka/src/config"
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

	// message := TransactionMessage{
	// 	Sender_id:   c.Locals("user").(database.User).UserID,
	// 	Receiver_id: payload.Receiver_id,
	// 	Amount:      payload.Amount,
	// 	Status:      "PENDING",
	// }

	// jsonMessage, err := json.Marshal(message)
	// if err != nil {
	// 	fmt.Printf("Failed to marshal JSON: %v\n", err)
	// }

	// err = producer.PushTransactionToQueue("wallet.pending", jsonMessage)
	// if err != nil {
	// 	fmt.Printf("Failed to push transaction to Kafka: %v\n", err)
	// }

	return c.JSON(fiber.Map{
		"message": "Transaction Created",
	})
}

func GetBalance(c *fiber.Ctx) error {
	db := config.DB.Db

	user, err := db.Exec("SELECT * FROM users WHERE id = $1", "id")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}
