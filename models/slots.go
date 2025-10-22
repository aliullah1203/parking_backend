package models

import "parking_management_system_backend/config"

// ResetExpiredSlots sets slots as free if their end_time has passed
func ResetExpiredSlots() error {
	_, err := config.DB.Exec(`
		UPDATE slots
		SET user_id = NULL, start_time = NULL, end_time = NULL, notified = false
		WHERE end_time IS NOT NULL AND end_time <= NOW()
	`)
	return err
}
