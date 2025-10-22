package services

import (
	"context"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2api "google.golang.org/api/oauth2/v2"

	"github.com/joho/godotenv"
)

var GoogleOauthConfig *oauth2.Config

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading environment variables from system")
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Google OAuth client ID or secret is missing in .env")
	}

	GoogleOauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/api/oauth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

// GetGoogleLoginURL generates the Google OAuth login URL
func GetGoogleLoginURL(state string) string {
	return GoogleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// HandleGoogleCallback exchanges code and fetches user info
func HandleGoogleCallback(code string) (user *oauth2api.Userinfo, accessToken, refreshToken string, err error) {
	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to exchange code: %v", err)
	}

	client := GoogleOauthConfig.Client(context.Background(), token)
	service, err := oauth2api.New(client)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to create OAuth2 service: %v", err)
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to get user info: %v", err)
	}

	return userInfo, token.AccessToken, token.RefreshToken, nil
}
