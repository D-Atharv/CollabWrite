package auth

import (
	"context"
	"fmt"
	"net/http"

	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var jwtSecret = []byte("JWT_SECRET")

var GoogleOAuthConfig = &oauth2.Config{
	// TODO: fix env issue
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL: "http://localhost:8080/auth/google/callback",
	Scopes:      []string{"email", "profile"},
	Endpoint:    google.Endpoint,
}

func GoogleLogin(c *gin.Context) {
	url := GoogleOAuthConfig.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func GoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing auth code",
		})
		return
	}

	_, err := GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "OAuth exchange failed",
		})
		return
	}

	jwtToken, _ := GenerateJWT(1)
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"token": jwtToken,
	// })

	//change this later
	ctx.Header("Content-Type", "text/html")
	ctx.String(http.StatusOK, `
        <html>
        <head><title>Login Successful</title></head>
        <body>
            <h1>Successfully authenticated with Google</h1>
            <p>Your JWT Token:</p>
            <textarea style="width: 100%%; height: 150px;">`+jwtToken+`</textarea>
        </body>
        </html>
    `)

	fmt.Print("Successfully authenticated with Google")
}

func GenerateJWT(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(jwtSecret)
}
