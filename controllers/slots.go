package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"parking_management_system_backend/jobs"
	"parking_management_system_backend/models"
	"parking_management_system_backend/user"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Get all slots
func GetSlots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slots, err := models.GetAllSlots()
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(slots)
}

// BookSlotRequest is used to decode incoming JSON
// Frontend commonly sends times as ISO strings, so we parse them here.
type BookSlotRequest struct {
	ID        int     `json:"id"`                // slot id
	UserID    *string `json:"user_id,omitempty"` // nullable
	StartTime string  `json:"start_time"`        // ISO string
	EndTime   string  `json:"end_time"`          // ISO string
	Email     string  `json:"email,omitempty"`   // email string
}

// Book a slot
func BookSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req BookSlotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Invalid request body: %v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// parse times (expecting RFC3339 or RFC3339Nano)
	parse := func(s string) (time.Time, error) {
		if s == "" {
			return time.Time{}, nil
		}
		t, err := time.Parse(time.RFC3339, s)
		if err == nil {
			return t, nil
		}
		// try RFC3339Nano fallback
		return time.Parse(time.RFC3339Nano, s)
	}

	startT, err := parse(req.StartTime)
	if err != nil {
		http.Error(w, `{"error":"Invalid start_time format. Use ISO8601 (RFC3339)."}`, http.StatusBadRequest)
		return
	}
	endT, err := parse(req.EndTime)
	if err != nil {
		http.Error(w, `{"error":"Invalid end_time format. Use ISO8601 (RFC3339)."}`, http.StatusBadRequest)
		return
	}

	// Basic validation
	if startT.IsZero() || endT.IsZero() {
		http.Error(w, `{"error":"Start time and end time are required"}`, http.StatusBadRequest)
		return
	}
	if startT.After(endT) {
		http.Error(w, `{"error":"Start time cannot be after end time"}`, http.StatusBadRequest)
		return
	}

	// Build user.Slot with sql.Null* types
	var s user.Slot
	s.ID = req.ID

	if req.UserID != nil && *req.UserID != "" {
		s.UserID.Valid = true
		s.UserID.String = *req.UserID
	} else {
		s.UserID.Valid = false
		s.UserID.String = ""
	}

	if req.Email != "" {
		s.Email.Valid = true
		s.Email.String = req.Email
	} else {
		s.Email.Valid = false
		s.Email.String = ""
	}

	s.StartTime.Valid = true
	s.StartTime.Time = startT
	s.EndTime.Valid = true
	s.EndTime.Time = endT
	s.Notified = false // newly booked -> not notified

	// Save to database
	if err := models.BookSlot(&s); err != nil {
		// if the model returns a friendly error like "slot already booked", forward it
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

    // Optional: send a confirmation/reminder with minutes until end
    if s.Email.Valid && s.Email.String != "" {
        remaining := int(endT.Sub(time.Now()).Minutes())
        if remaining < 0 {
            remaining = 0
        }
        go func(email string, slotID int, minutes int) {
            _ = jobs.SendEmailReminder(email, slotID, minutes)
        }(s.Email.String, s.ID, remaining)
    }

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Slot booked successfully",
	})
}

// Release a slot
func ReleaseSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var payload struct{ ID int }
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"Invalid request body: %v"}`, err.Error()), http.StatusBadRequest)
		return
	}

	if err := models.ReleaseSlot(payload.ID); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Slot released successfully"})
}

// Notify slot reminder
func NotifySlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]
	slotID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid slot ID"}`, http.StatusBadRequest)
		return
	}

	slot, err := models.GetSlotByID(slotID)
	if err != nil {
		http.Error(w, `{"error":"Slot not found"}`, http.StatusNotFound)
		return
	}

	if slot.UserID.Valid && slot.Email.Valid {
		remaining := int(slot.EndTime.Time.Sub(time.Now()).Minutes())
		if remaining < 0 {
			remaining = 0
		}
		_ = jobs.SendEmailReminder(slot.Email.String, slot.ID, remaining)
		slot.Notified = true
		_ = models.UpdateSlot(slot)
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Reminder sent successfully"})
}
