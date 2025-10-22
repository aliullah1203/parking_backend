package jobs

// import (
// 	"fmt"
// 	"parking_management_system_backend/config"
// 	"parking_management_system_backend/user"
// 	"time"
// )

// // SendEmailReminder is a thin helper to format a reminder and send it via SMTP
// func SendEmailReminder(email string, slotID int, remainingMinutes int) error {
// 	subject := "Parking Slot Reminder"
// 	body := fmt.Sprintf("Reminder: Your slot #%d ends in %d minutes.", slotID, remainingMinutes)
// 	return SendEmail(email, subject, body)
// }

// // SendSlotReminders finds slots ending in next 30 minutes and notifies once
// func SendSlotReminders() {
// 	now := time.Now()
// 	rows, err := config.DB.Queryx(`SELECT * FROM slots WHERE notified=false AND end_time > $1 AND end_time <= $2`, now, now.Add(30*time.Minute))
// 	if err != nil {
// 		fmt.Println("Failed to fetch slots:", err)
// 		return
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var s user.Slot
// 		if err := rows.StructScan(&s); err != nil {
// 			fmt.Println("Failed to scan slot:", err)
// 			continue
// 		}

// 		if s.Email.Valid {
// 			remaining := int(s.EndTime.Time.Sub(now).Minutes())
// 			if remaining < 0 {
// 				remaining = 0
// 			}
// 			if err := SendEmailReminder(s.Email.String, s.ID, remaining); err != nil {
// 				fmt.Println("Failed to send email:", err)
// 			} else {
// 				// Mark as notified to avoid repeats
// 				if _, err := config.DB.Exec("UPDATE slots SET notified=true WHERE id=$1", s.ID); err != nil {
// 					fmt.Println("Failed to mark slot notified:", err)
// 				}
// 			}
// 		}
// 	}
// }
