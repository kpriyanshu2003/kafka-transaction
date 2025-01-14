package models

import "github.com/google/uuid"

type Transaction struct {
	Sender_id   uuid.UUID `json:"sender_id"`
	Receiver_id uuid.UUID `json:"receiver_id"`
	Amount      string    `json:"amount"`
	Status      string    `json:"status"`
	Type        string    `json:"type"` // Debit, Credit
}
