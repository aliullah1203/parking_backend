package controllers

// You can add the logout logic below, e.g.,
import (
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	// Example: just return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Logged out successfully"}`))
}
