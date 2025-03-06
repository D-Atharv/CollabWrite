package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"server/internal/models"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var jwtSecret = []byte("JWT_SECRET")

var GoogleOAuthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	Scopes:       []string{"email", "profile"},
	Endpoint:     google.Endpoint,
}
func FetchGoogleUser(token *oauth2.Token) (*models.User, error) {
	client := GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("error fetching user info from Google API: %w", err)
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("error decoding user info from Google API: %w", err)
	}

	user := &models.User{
		Provider:   "google",
		ProviderID: userInfo.ID,
		Email:      userInfo.Email,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return user, nil
}

func FindOrCreateUser(db *gorm.DB, googleUser *models.User) (*models.User, error) {
	var user models.User
	if err := db.Where("provider_id = ?", googleUser.ProviderID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(googleUser).Error; err != nil {
				return nil, fmt.Errorf("error creating user in DB: %w", err)
			}
			return googleUser, nil
		}
		return nil, fmt.Errorf("error finding user in DB: %w", err)
	}
	return &user, nil
}

