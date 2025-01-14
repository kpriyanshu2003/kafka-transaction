package models

import "github.com/google/uuid"

type Transaction struct {
	TransactionID uuid.UUID
	SendID        uuid.UUID
	RecvID        uuid.UUID
	Amt           string
	Status        string
	Type          string
}
