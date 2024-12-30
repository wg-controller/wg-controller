package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lampy255/net-tbm/db"
	"github.com/lampy255/net-tbm/types"
)

func StartAPI() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.New()

	// Create router groups
	public := router.Group("/api/v1")
	private := router.Group("/api/v1")

	// Middleware
	private.Use(AuthMiddleware)

	// Public Endpoints
	public.GET("/health", GET_Health)
	public.POST("/prelogin", POST_PreLogin)
	public.POST("/login", POST_Login)

	// Private Endpoints
	private.GET("/peers", GET_Peers)
	private.GET("/peers/:uuid", GET_Peer)
	private.PUT("/peers/:uuid", PUT_Peer)
	private.PATCH("/peers/:uuid", PATCH_Peer)
	private.DELETE("/peers/:uuid", DELETE_Peer)

	private.GET("/accounts", GET_Accounts)
	private.PUT("/accounts/:email", PUT_Account)
	private.PATCH("/accounts/:email", PATCH_Account)
	private.DELETE("/accounts/:email", DELETE_Account)

	private.GET("/apikeys", GET_APIKeys)
	private.PUT("/apikeys/:uuid", PUT_APIKey)
	private.DELETE("/apikeys/:uuid", DELETE_APIKey)

	private.GET("/poll", GET_LongPoll)

	// Start Listening
	log.Println("Starting API on port " + ENV.API_PORT)
	err := router.Run(":" + ENV.API_PORT)
	if err != nil {
		log.Fatal("Error starting API:", err)
	}
}

func GET_Health(c *gin.Context) {
	_, err := GetWireguard()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func GET_Peers(c *gin.Context) {
	peers, err := db.GetPeers()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	var extendedPeers []types.PeerExtended
	for _, peer := range peers {
		extendedPeer, err := GetWireguardPeer(peer)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		extendedPeers = append(extendedPeers, extendedPeer)
	}

	c.JSON(200, extendedPeers)
}

func GET_Peer(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(400, gin.H{
			"error": "uuid is required",
		})
		return
	}

	peer, err := db.GetPeer(uuid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	extendedPeer, err := GetWireguardPeer(peer)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, extendedPeer)
}

func PUT_Peer(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(400, gin.H{
			"error": "uuid is required",
		})
		return
	}

	var peer types.Peer
	peer.UUID = uuid
	err := c.BindJSON(&peer)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.InsertPeer(peer)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func PATCH_Peer(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(400, gin.H{
			"error": "uuid is required",
		})
		return
	}

	var peer types.Peer
	peer.UUID = uuid
	err := c.BindJSON(&peer)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdatePeer(peer)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func DELETE_Peer(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(400, gin.H{
			"error": "uuid is required",
		})
		return
	}

	err := db.DeletePeer(uuid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func GET_Accounts(c *gin.Context) {
	accounts, err := db.GetAccounts()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, accounts)
}

func PUT_Account(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(400, gin.H{
			"error": "email is required",
		})
		return
	}

	// Parse the account request body
	var account types.UserAccountWithPass
	err := c.BindJSON(&account)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	account.Email = email

	// Create salt
	salt, err := NewSalt()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash password
	hash, err := HashString(account.Password, salt)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Insert account
	err = db.InsertAccount(email, account.Role, hash, salt)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func PATCH_Account(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(400, gin.H{
			"error": "email is required",
		})
		return
	}

	// Parse the account request body
	var account types.UserAccountWithPass
	err := c.BindJSON(&account)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	account.Email = email

	// Create salt
	salt, err := NewSalt()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash password
	hash, err := HashString(account.Password, salt)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Update account
	err = db.UpdateAccountPasswordHash(email, hash, salt)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func DELETE_Account(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(400, gin.H{
			"error": "email is required",
		})
		return
	}

	err := db.DeleteAccount(email)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func GET_APIKeys(c *gin.Context) {
	apiKeys, err := db.GetApiKeys()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, apiKeys)
}

func PUT_APIKey(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(400, gin.H{
			"error": "uuid is required",
		})
		return
	}

	// Parse the api key request body
	var apiKey types.APIKey
	err := c.BindJSON(&apiKey)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	apiKey.UUID = uuid

	// Insert api key
	err = db.InsertApiKey(apiKey)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func DELETE_APIKey(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(400, gin.H{
			"error": "uuid is required",
		})
		return
	}

	err := db.DeleteApiKey(uuid)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}
