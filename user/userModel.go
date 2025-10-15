package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	Name               string     `json:"name" db:"name"`
	Email              string     `json:"email" db:"email"`
	Phone              *string    `json:"phone,omitempty" db:"phone"`
	Address            *string    `json:"address,omitempty" db:"address"`
	Role               string     `json:"role" db:"role"`
	Status             string     `json:"status" db:"status"`
	SubscriptionStatus string     `json:"subscription_status" db:"subscription_status"`
	Password           string     `json:"password,omitempty" db:"password"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
