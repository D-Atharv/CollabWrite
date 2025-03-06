package auth_handlers

import (
	"net/http"
	"server/internal/auth"
	"server/internal/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func GoogleLogin(c *gin.Context) {
	url := auth.GoogleOAuthConfig.AuthCodeURL("random-state-string", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing auth code"})
		return
	}

	token, err := auth.GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "OAuth exchange failed"})
		return
	}

	googleUser, err := auth.FetchGoogleUser(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := auth.FindOrCreateUser(db.DB, googleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	jwtToken, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": jwtToken, "message": "Successfully authenticated"})
}
