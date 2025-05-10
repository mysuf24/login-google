package controller

import (
	"encoding/json"
	"log"
	"login-google/config"
	"login-google/model"
	"login-google/repository"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {
	state := generateState()
	url := config.GoogleOAuthConfig.AuthCodeURL(state)

	log.Printf("[GoogleLogin] Redirecting to : %s", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func generateState() string {
	rand.Seed(time.Now().UnixNano())
	letter := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 16)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		log.Println("[GoogleCallback] No code in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in query"})
		return
	}

	token, err := config.GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		log.Printf("[GoogleCallback] Token exchange error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := config.GoogleOAuthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("[GoogleCallback] Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userData struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		log.Printf("[GoogleCallback] Failed to parse user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	userRepo := repository.UserRepository{DB: config.DB}
	user, err := userRepo.FindByEmail(userData.Email)
	if err != nil {
		log.Println("[GoogleCallback] User not found, creating new user")
		user = &model.User{
			Email:  userData.Email,
			Name:   userData.Name,
			Avatar: userData.Picture,
		}
		err = userRepo.Create(user)
		if err != nil {
			log.Printf("[GoogleCallback] Failed to create user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Generate JWT
	jwtToken, err := config.GenerateJWT(user.Email)
	if err != nil {
		log.Printf("[GoogleCallback] Failed to generate JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   jwtToken,
		"user": gin.H{
			"name":   user.Name,
			"email":  user.Email,
			"avatar": user.Avatar,
		},
	})
}
