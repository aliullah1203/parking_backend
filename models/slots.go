package models

import (
	"database/sql"
	"errors"
	"parking_management_system_backend/config"
	"parking_management_system_backend/user"
)

// GetAllSlots fetches all slots
func GetAllSlots() ([]user.Slot, error) {
	var slots []user.Slot
	err := config.DB.Select(&slots, "SELECT * FROM slots ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	return slots, nil
}

// BookSlot books a slot only if it is free (start_time AND end_time are NULL)
func BookSlot(s *user.Slot) error {
	query := `
		UPDATE slots
		SET user_id=$1, start_time=$2, end_time=$3, email=$4, notified=false
		WHERE id=$5 AND start_time IS NULL AND end_time IS NULL
		RETURNING id;
	`
	return config.DB.QueryRow(
		query,
		s.UserID,
		s.StartTime,
		s.EndTime,
		s.Email,
		s.ID,
	).Scan(&s.ID)
}

// ReleaseSlot clears a slot
func ReleaseSlot(id int) error {
	_, err := config.DB.Exec(`
		UPDATE slots
		SET user_id=NULL, start_time=NULL, end_time=NULL, email=NULL, notified=false
		WHERE id=$1
	`, id)
	return err
}

// GetSlotByID fetches a slot by its ID
func GetSlotByID(id int) (*user.Slot, error) {
	var s user.Slot
	err := config.DB.Get(&s, "SELECT * FROM slots WHERE id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("slot not found")
		}
		return nil, err
	}
	return &s, nil
}

// UpdateSlot updates a slot's notified status
func UpdateSlot(s *user.Slot) error {
	_, err := config.DB.Exec("UPDATE slots SET notified=$1 WHERE id=$2", s.Notified, s.ID)
	return err
}
