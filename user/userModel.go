package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID `db:"id"`
	Name               string    `db:"name"`
	Email              string    `db:"email"`
	Phone              string    `db:"phone"`
	Role               string    `db:"role"`
	Status             string    `db:"status"`
	SubscriptionStatus string    `db:"subscription_status"`
	Password           string    `db:"password"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}
