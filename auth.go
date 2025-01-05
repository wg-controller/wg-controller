package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lampy255/net-tbm/db"
	"github.com/lampy255/net-tbm/types"
)

const MaxFailedAttempts = 5

func GenerateDeterministicHash(input []byte, salt []byte) (hash []byte, err error) {
	// Check for empty input
	if len(input) == 0 {
		return []byte{}, errors.New("input cannot be empty")
	}

	// Combine the input and salt
	joined := append(input, salt...)

	// Hash the input using SHA256
	h := sha256.New()
	_, err = h.Write(joined)
	if err != nil {
		return []byte{}, err
	}

	return h.Sum(nil), nil
}

func NewSalt() (salt []byte, err error) {
	return GenerateRandomBytes(16)
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.StdEncoding.EncodeToString(b), err
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
		// Decode Base64
		tokenBytes, err := base64.URLEncoding.DecodeString(session)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
			return
		}

		// Hash the session token
		tokenHash, err := GenerateDeterministicHash(tokenBytes, []byte{})
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
			return
		}

		// Check for the hashed session token in the DB
		expires, email, err := db.GetSession(tokenHash)
		if err != nil {
			c.AbortWithStatus(403)
			log.Println(err)
			log.Println("Invalid sessionId from IP:", c.ClientIP())
			return
		}

		// Check if the session is expired
		if expires < time.Now().UnixMilli() {
			c.AbortWithStatus(403)
			return
		}

		// Update last active time for user
		err = db.UpdateAccountLastActive(email, time.Now().UnixMilli())
		if err != nil {
			log.Println(err)
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

		// Decode Base64
		tokenBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(500)
			return
		}

		// Hash the api key
		hash, err := GenerateDeterministicHash(tokenBytes, []byte{})
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
		c.JSON(401, gin.H{
			"error": "not logged in",
		})
		return
	}

	// Decode Base64
	sessionBytes, err := base64.URLEncoding.DecodeString(session)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Hash the session token
	hash, err := GenerateDeterministicHash(sessionBytes, []byte{})
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Check for the hashed session token in the DB
	expires, email, err := db.GetSession(hash)
	if err != nil {
		log.Println(err)
		c.JSON(401, gin.H{
			"error": "invalid session",
		})
		return
	}

	// Check if the session token is expired
	if expires < time.Now().UnixMilli() {
		c.JSON(401, gin.H{
			"error": "session expired",
		})
		return
	}

	// Update last active time for user
	err = db.UpdateAccountLastActive(email, time.Now().UnixMilli())
	if err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{
		"status": "ok",
		"email":  email,
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

		c.JSON(401, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// Check if the user is suspended
	if account.FailedAttempts >= MaxFailedAttempts {
		log.Println("User is suspended:", login.Email, "from IP:", c.ClientIP())
		c.JSON(401, gin.H{
			"error": "account suspended",
		})
		return
	}

	// Get the stored password hash and salt
	storedHash, salt, err := db.GetAccountPasswordHash(login.Email)
	if err != nil {
		log.Println("Error getting password hash for user:", login.Email, "from IP:", c.ClientIP())
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Hash the input password
	testHash, err := GenerateDeterministicHash([]byte(login.Password), salt)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Compare the stored hash with the input hash
	for i := 0; i < len(storedHash); i++ {
		if storedHash[i] != testHash[i] {
			log.Println("Invalid credentials for user:", login.Email, "from IP:", c.ClientIP())
			c.JSON(401, gin.H{
				"error": "invalid email or password",
			})

			// Increment the failed attempts
			err := db.IncrementAccountFailedAttempts(login.Email)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}

	// Generate a session token
	tokenBytes, err := GenerateRandomBytes(32)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Hash the session token
	tokenHash, err := GenerateDeterministicHash(tokenBytes, []byte{})
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Store the hashed session token
	err = db.CreateSession(tokenHash, login.Email, time.Now().Add(time.Hour*12).UnixMilli())
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Base64 encode the session token
	tokenBase64 := base64.URLEncoding.EncodeToString(tokenBytes)

	// Set cookie
	c.SetCookie("sessionId", tokenBase64, 0, "", "", true, true)
	log.Println("User logged in:", login.Email, "from IP:", c.ClientIP())
	c.JSON(200, gin.H{
		"status": "ok",
		"email":  login.Email,
	})
}

func POST_Logout(c *gin.Context) {
	// Check for session cookie
	session, err := c.Cookie("sessionId")
	if err != nil {
		c.JSON(401, gin.H{
			"error": "not logged in",
		})
		return
	}

	// Decode Base64
	sessionBytes, err := base64.URLEncoding.DecodeString(session)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Hash the session token
	hash, err := GenerateDeterministicHash(sessionBytes, []byte{})
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Check that the session token exists
	_, email, err := db.GetSession(hash)
	if err != nil {
		log.Println(err)
		c.JSON(401, gin.H{
			"error": "invalid session",
		})
		return
	}

	// Delete the session from the DB
	err = db.DeleteSession(hash)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}

	// Delete the session cookie
	c.SetCookie("sessionId", "", -1, "", "", true, true)
	log.Println("User logged out:", email, "from IP:", c.ClientIP())
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
