package authentication

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Password string `json:"password" binding:"required"`
}

// current security key for the cookie
var key string

// Authentication middleware
func AuthRequired(c *gin.Context) {
	if CheckAuthCookie(c) {
		c.Next()
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		c.Abort()
	}
}

// login function, checks password and sets cookie accordingly
func authLogin(c *gin.Context) {
	password := os.Getenv("APIPASSWORD")

	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No password is set"})
			return
	} else {
		var loginReq loginRequest

		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing POST header: " + err.Error()})
			return
		}
		
		// if password is correct
		if password == loginReq.Password {
			c.SetSameSite(http.SameSiteStrictMode)
			c.SetCookie("auth", key, 3600000, "/", "localhost", true, true)
			c.JSON(http.StatusOK, gin.H{"success": "Login succeeded"})
			log.Println("User logged in successfully")
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong password"})
			return
		}
	}
}

// reports the current authorization status
func authStatus(c *gin.Context) {
	// check if correct cookie is set
	if CheckAuthCookie(c) {
		c.String(http.StatusOK, "auth OK")
	} else {
		c.String(http.StatusUnauthorized, "auth FAILED")
	}
}

func CheckAuthCookie(c *gin.Context) bool {
	cookie, err := c.Cookie("auth")

	// check if correct cookie is set
	if err == nil && cookie == key {
		return true
	} else {
		return false
	}
}

func InitializeAuth(router *gin.Engine) {
	{	
		// generate random cookie value
		rand.Seed(time.Now().UnixNano())
		randInt := rand.Intn(10000000) + 1000000
		hash := md5.Sum([]byte(strconv.Itoa(randInt)))
		key = hex.EncodeToString(hash[:])
	}

	router.GET("/api/auth/status", authStatus)
	router.POST("/api/auth/login", authLogin)
}