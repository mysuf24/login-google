package controller

import (
	"io"
	"log"
	"login-google/config"
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
		log.Println("[GoogleCallback] no code in query")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in query"})
		return
	}

	token, err := config.GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		log.Printf("[GoogleCallback] Token exchange error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal to exchange token"})
		return
	}

	client := config.GoogleOAuthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Printf("[GoogleCallback] Gagal untuk mengambil info user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal untuk mengambil info user"})
		return
	}

	userInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[GoogleCallback] Gagal untuk membaca info user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal untuk membaca info user"})
		return
	}

	log.Printf("'[GoogleCallback] User Info: %s", userInfo)
	c.Data(http.StatusOK, "application/json", userInfo)
}
