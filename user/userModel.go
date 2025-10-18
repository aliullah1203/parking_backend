package models

import (
	"time"
)

type User struct {
	ID                 string    `db:"id" json:"id"`
	Name               string    `db:"name" json:"name"`
	Email              string    `db:"email" json:"email"`
	Phone              string    `db:"phone" json:"phone"`
	Role               string    `db:"role" json:"role"`
	Status             string    `db:"status" json:"status"`
	SubscriptionStatus string    `db:"subscription_status" json:"subscriptionStatus"`
	Password           string    `db:"password" json:"-"`
	License            string    `db:"license"`
	NID                string    `db:"nid"`
	Picture            string    `db:"picture"`
	CreatedAt          time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time `db:"updated_at" json:"updatedAt"`
}
