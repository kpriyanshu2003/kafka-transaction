package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	UserID  uuid.UUID
	Email   string
	Balance sql.NullString
}
