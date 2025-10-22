package controllers

// import (
// 	"encoding/json"
// 	"net/http"
// 	"parking_management_system_backend/helpers"
// 	"parking_management_system_backend/jobs"
// 	"parking_management_system_backend/models"
// 	"strconv"
// 	"time"
// )

// // GetSlots returns all slots and resets expired slots
// func GetSlots(w http.ResponseWriter, r *http.Request) {
// 	// Reset expired slots first
// 	if err := models.ResetExpiredSlots(); err != nil {
// 		http.Error(w, "Failed to reset expired slots", http.StatusInternalServerError)
// 		return
// 	}

// 	slots, err := models.GetAllSlots()
// 	if err != nil {
// 		http.Error(w, "Failed to fetch slots", http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(slots)
// }

// // BookSlotHandler handles booking a slot
// func BookSlot(w http.ResponseWriter, r *http.Request) {
// 	claims := helpers.GetClaimsFromContext(r)
// 	if claims == nil {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	var payload struct {
// 		ID        int       `json:"id"`
// 		StartTime time.Time `json:"start_time"`
// 		EndTime   time.Time `json:"end_time"`
// 		Email     string    `json:"email"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		http.Error(w, "Invalid payload", http.StatusBadRequest)
// 		return
// 	}

// 	// Fetch slot from DB
// 	slot, err := models.GetSlotByID(payload.ID)
// 	if err != nil {
// 		http.Error(w, "Slot not found", http.StatusNotFound)
// 		return
// 	}

// 	// Check availability
// 	if slot.UserID.Valid && slot.EndTime.Valid && time.Now().Before(slot.EndTime.Time) {
// 		http.Error(w, "Slot already booked", http.StatusBadRequest)
// 		return
// 	}

// 	// Update slot fields
// 	slot.UserID.String = claims.UserID
// 	slot.UserID.Valid = true
// 	slot.StartTime.Time = payload.StartTime
// 	slot.StartTime.Valid = true
// 	slot.EndTime.Time = payload.EndTime
// 	slot.EndTime.Valid = true
// 	slot.Email.String = payload.Email
// 	slot.Email.Valid = true
// 	slot.Notified = false

// 	if err := models.BookSlot(slot); err != nil {
// 		http.Error(w, "Failed to book slot", http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]string{"message": "Slot booked successfully"})
// }

// // ReleaseSlotHandler releases a slot manually
// func ReleaseSlot(w http.ResponseWriter, r *http.Request) {
// 	var payload struct {
// 		ID int `json:"id"`
// 	}
// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		http.Error(w, "Invalid payload", http.StatusBadRequest)
// 		return
// 	}

// 	if err := models.ReleaseSlot(payload.ID); err != nil {
// 		http.Error(w, "Failed to release slot", http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]string{"message": "Slot released successfully"})
// }

// // NotifySlotHandler sends a notification for a slot
// func NotifySlot(w http.ResponseWriter, r *http.Request) {
// 	vars := r.URL.Query()
// 	idStr := vars.Get("id")
// 	if idStr == "" {
// 		http.Error(w, "Missing slot ID", http.StatusBadRequest)
// 		return
// 	}

// 	slotID, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid slot ID", http.StatusBadRequest)
// 		return
// 	}

// 	slot, err := models.GetSlotByID(slotID)
// 	if err != nil {
// 		http.Error(w, "Slot not found", http.StatusNotFound)
// 		return
// 	}

// 	if slot.Email.Valid {
// 		_ = jobs.SendEmailReminder(slot.Email.String, slot.ID, int(slot.EndTime.Time.Sub(time.Now()).Minutes()))
// 	}

// 	json.NewEncoder(w).Encode(map[string]string{"message": "Notification sent"})
// }
