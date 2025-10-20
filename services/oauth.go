package services

import (
	"errors"
)

type OAuthUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func GetGoogleLoginURL(state string) string {
	// Placeholder; replace with real Google OAuth login URL builder
	return "/"
}

func HandleGoogleCallback(code string) (*OAuthUser, string, string, error) {
	// Placeholder; replace with real Google OAuth callback handler
	return nil, "", "", errors.New("OAuth not configured")
}
