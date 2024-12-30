package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lampy255/net-tbm/db"
	"github.com/lampy255/net-tbm/types"
	"golang.org/x/crypto/bcrypt"
)

const MaxFailedAttempts = 5

// Hashes a password with a salt
func HashPassword(password string, salt string) (hash string, err error) {
	// Check for empty password
	if password == "" {
		return "", errors.New("empty password")
	}

	// Convert to byte slices
	passwordBytes := []byte(password)
	saltBytes := []byte(salt)
	combined := append(passwordBytes, saltBytes...)

	// Hash the password
	output, cryptErr := bcrypt.GenerateFromPassword(combined, bcrypt.DefaultCost)
	if cryptErr != nil {
		return "", cryptErr
	}

	// Convert to string
	hash = string(output)

	return hash, err
}

func NewSalt() (salt string, err error) {
	return GenerateRandomString(16)
}

func NewSessionToken() (token string, err error) {
	return GenerateRandomString(32)
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func AuthMiddleware(c *gin.Context) {
	// Check for sessionId cookie
	bearer := ""
	session, err := c.Cookie("sessionId")
	if err != nil {
		// Check for Authorization header
		bearer = c.GetHeader("Authorization")
		if bearer == "" {
			c.AbortWithStatus(403)
			log.Println("No Authorization header or sessionId cookie from IP:", c.ClientIP())
			return
		}
	}

	// If sessionId cookie is present, check the session
	if session != "" {
		// Hash the session token
		hash, err := HashPassword(session, "")
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
			return
		}

		// Check for the hashed session token in the DB
		expires, err := db.GetSession(hash)
		if err != nil {
			c.AbortWithStatus(403)
			log.Println("Invalid sessionId from IP:", c.ClientIP())
			return
		}

		// Check if the session is expired
		if expires < time.Now().UnixMilli() {
			c.AbortWithStatus(403)
			return
		}

		c.Next()
		return
	}

	// If Authorization header is present, check the token
	if bearer != "" {
		if len(bearer) < 7 {
			c.AbortWithStatus(403)
			log.Println("Invalid Authorization header from IP:", c.ClientIP())
			return
		}

		token := bearer[7:]
		if token == "" {
			c.AbortWithStatus(403)
			log.Println("Invalid Authorization header from IP:", c.ClientIP())
			return
		}

		// Hash the api key
		hash, err := HashPassword(token, "")
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
			return
		}

		// Check for the api key in the DB
		expires, err := db.GetApiKey(hash)
		if err != nil {
			c.AbortWithStatus(403)
			log.Println("Invalid token from IP:", c.ClientIP())
			return
		}

		// Check if the api key is expired
		if expires < time.Now().UnixMilli() {
			c.AbortWithStatus(403)
			return
		}

		c.Next()
		return
	}

	// Default to 403
	c.AbortWithStatus(403)
}

func POST_PreLogin(c *gin.Context) {
	// Check for session cookie
	session, err := c.Cookie("sessionId")
	if err != nil {
		c.JSON(403, gin.H{
			"error": "not logged in",
		})
	}

	// Hash the session token
	hash, err := HashPassword(session, "")
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Check for the hashed session token in the DB
	expires, err := db.GetSession(hash)
	if err != nil {
		c.JSON(403, gin.H{
			"error": "invalid session",
		})
	}

	// Check if the session token is expired
	if expires < time.Now().UnixMilli() {
		c.JSON(403, gin.H{
			"error": "session expired",
		})
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

// Gin handler for login requests
func POST_Login(c *gin.Context) {
	// Parse the login request body
	var login types.LoginBody
	err := c.BindJSON(&login)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check for empty fields
	if login.Email == "" || login.Password == "" {
		log.Println("Malformed login request from IP:" + c.ClientIP())
		c.JSON(400, gin.H{
			"error": "malformed request",
		})
		return
	}

	// Check if the user exists
	account, err := db.GetAccount(login.Email)
	if err != nil {
		log.Println("Error getting account", login.Email, "from IP:", c.ClientIP())
		log.Println(err)

		c.JSON(403, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// Check if the user is suspended
	if account.FailedAttempts >= MaxFailedAttempts {
		log.Println("User is suspended:", login.Email, "from IP:", c.ClientIP())
		c.JSON(403, gin.H{
			"error": "account suspended",
		})
		return
	}

	// Get the stored password hash and salt
	hash, salt, err := db.GetAccountPasswordHash(login.Email)
	if err != nil {
		log.Println("Error getting password hash for user:", login.Email, "from IP:", c.ClientIP())
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Compare the password hashes
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(login.Password+salt))
	if err != nil {
		log.Println("Invalid credentials for user:", login.Email, "from IP:", c.ClientIP())

		// Increment the failed attempts
		err := db.IncrementAccountFailedAttempts(login.Email)
		if err != nil {
			log.Println(err)
		}

		c.JSON(403, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// Reset the failed attempts
	err = db.ResetAccountFailedAttempts(login.Email)
	if err != nil {
		log.Println(err)
	}

	// Generate a session token
	token, err := NewSessionToken()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Hash the session token
	hash, err = HashPassword(token, salt)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Store the session token
	err = db.CreateSession(hash, login.Email, time.Now().Add(time.Hour*12).UnixMilli())
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Set cookie
	c.SetCookie("sessionId", token, 0, "", "", false, true)
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
