package main

import (
	"encoding/base64"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wg-controller/wg-controller/db"
	"github.com/wg-controller/wg-controller/types"
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
	public.POST("/logout", POST_Logout)

	// Private Endpoints
	private.GET("/peers", GET_Peers)
	private.GET("/peers/:uuid", GET_Peer)
	private.PUT("/peers/:uuid", PUT_Peer)
	private.PATCH("/peers/:uuid", PATCH_Peer)
	private.DELETE("/peers/:uuid", DELETE_Peer)
	private.GET("/peers/init", GET_InitPeer)

	private.GET("/accounts", GET_Accounts)
	private.PUT("/accounts/:email", PUT_Account)
	private.PATCH("/accounts/:email", PATCH_Account)
	private.PATCH("/accounts/:email/password", PATCH_AccountPassword)
	private.DELETE("/accounts/:email", DELETE_Account)

	private.GET("/apikeys", GET_APIKeys)
	private.PUT("/apikeys/:uuid", PUT_APIKey)
	private.PATCH("/apikeys/:uuid", PATCH_APIKey)
	private.DELETE("/apikeys/:uuid", DELETE_APIKey)
	private.GET("/apikeys/init", GET_InitAPIKey)

	private.GET("/serverinfo", GET_ServerInfo)

	private.GET("/poll", GET_LongPoll)

	// Start Listening
	log.Println("Starting API on port " + ENV.API_PORT)
	err := router.Run(":" + ENV.API_PORT)
	if err != nil {
		log.Fatal("Error starting API:", err)
	}
}

func GET_Health(c *gin.Context) {
	_, err := wg.Device(ENV.WG_INTERFACE)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	var extendedPeers []types.Peer
	for _, peer := range peers {
		if peer.Enabled {
			extendedPeer, err := GetWireguardPeer(peer)
			if err != nil {
				log.Println(err)
				continue // Skip this peer
			}
			extendedPeers = append(extendedPeers, extendedPeer)
		} else {
			extendedPeers = append(extendedPeers, peer)
		}
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
		log.Println(err)
		c.Status(404)
		return
	}

	if peer.Enabled {
		peer, err = GetWireguardPeer(peer)
		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(200, peer)
}

func PUT_Peer(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.JSON(400, gin.H{
			"error": "uuid is required",
		})
		return
	}

	// Parse the peer request body
	var peer types.Peer
	peer.UUID = uuid
	err := c.BindJSON(&peer)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Insert peer into database
	err = db.InsertPeer(peer)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Resync wireguard configuration
	err = SyncWireguardConfiguration()
	if err != nil {
		log.Println(err)
	}

	// Resync peers DNS entries
	err = SyncPeersDNS(true)
	if err != nil {
		log.Println(err)
	}

	// Resync routing table
	err = SyncRoutingTable()
	if err != nil {
		log.Println(err)
	}

	// Push config to peer
	PushPeerConfig(peer)

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
		log.Println(err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.UpdatePeer(peer)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Resync wireguard configuration
	err = SyncWireguardConfiguration()
	if err != nil {
		log.Println(err)
	}

	// Resync peers DNS entries
	err = SyncPeersDNS(true)
	if err != nil {
		log.Println(err)
	}

	// Resync routing table
	err = SyncRoutingTable()
	if err != nil {
		log.Println(err)
	}

	// Push config to peer
	PushPeerConfig(peer)

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

	// Get peer
	peer, err := db.GetPeer(uuid)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete peer from database
	err = db.DeletePeer(uuid)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Prune wireguard configuration
	err = PruneWireguardPeers([]string{peer.PublicKey})
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	}

	// Resync peers DNS entries
	err = SyncPeersDNS(true)
	if err != nil {
		log.Println(err)
	}

	// Resync routing table
	err = SyncRoutingTable()
	if err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func GET_InitPeer(c *gin.Context) {
	InitPeer := types.PeerInit{}

	// Generate UUID
	InitPeer.UUID = uuid.New().String()

	// Generate private key
	privKey, err := NewWireguardPrivateKey()
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}
	InitPeer.PrivateKey = privKey

	// Generate public key
	pubKey, err := GetWireguardPublicKey(privKey)
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}
	InitPeer.PublicKey = pubKey

	// Generate pre-shared key
	preSharedKey, err := NewWireguardPreSharedKey()
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}
	InitPeer.PreSharedKey = preSharedKey

	// Get current peers to generate a unique localTunAddress
	peers, err := db.GetPeers()
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}

	// Extract used addresses
	var usedAddresses []string
	for _, peer := range peers {
		usedAddresses = append(usedAddresses, peer.RemoteTunAddress)
	}

	// Generate a unique address
	address, err := GetUniqueAddress(usedAddresses, ENV.SERVER_CIDR)
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}
	InitPeer.RemoteTunAddress = address

	InitPeer.ServerCIDR = ENV.SERVER_CIDR

	c.JSON(200, InitPeer)
}

func GET_Accounts(c *gin.Context) {
	accounts, err := db.GetAccounts()
	if err != nil {
		log.Println(err)
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
		log.Println(err)
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
	hash, err := GenerateDeterministicHash([]byte(account.Password), salt)
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
		log.Println(err)
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
	var account types.UserAccount
	err := c.BindJSON(&account)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	account.Email = email

	// Update account
	err = db.UpdateAccount(account)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func PATCH_AccountPassword(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(400, gin.H{
			"error": "email is required",
		})
		return
	}

	// Parse the password request body
	var password types.Password
	err := c.BindJSON(&password)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

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
	hash, err := GenerateDeterministicHash([]byte(password.Password), salt)
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
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete all sessions for this account
	err = db.DeleteUserSessions(email)
	if err != nil {
		log.Println(err)
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

	// Check that this is not the master admin account
	if email == ENV.ADMIN_EMAIL {
		c.JSON(400, gin.H{
			"error": "cannot delete master admin account",
		})
		return
	}

	// Delete account
	err := db.DeleteAccount(email)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Delete all sessions for this account
	err = db.DeleteUserSessions(email)
	if err != nil {
		log.Println(err)
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func GET_APIKeys(c *gin.Context) {
	apiKeys, err := db.GetApiKeys()
	if err != nil {
		log.Println(err)
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
	var apiKey types.APIKeyWithToken
	err := c.BindJSON(&apiKey)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	apiKey.UUID = uuid

	// Decode token from base64
	tokenBytes, err := base64.URLEncoding.DecodeString(apiKey.Token)
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}

	// Generate hash
	hash, err := GenerateDeterministicHash(tokenBytes, []byte{})
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}

	apiKeyWithoutToken := types.APIKey{
		UUID:              apiKey.UUID,
		Name:              apiKey.Name,
		ExpiresUnixMillis: apiKey.ExpiresUnixMillis,
		Attributes:        apiKey.Attributes,
	}

	// Insert api key
	err = db.InsertApiKey(apiKeyWithoutToken, hash)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func PATCH_APIKey(c *gin.Context) {
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
		log.Println(err)
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	apiKey.UUID = uuid

	// Update api key
	err = db.UpdateApiKey(apiKey)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func GET_InitAPIKey(c *gin.Context) {
	InitAPIKey := types.APIKeyInit{}

	// Generate UUID
	InitAPIKey.UUID = uuid.New().String()

	// Generate random token
	tokenBytes, err := GenerateRandomBytes(32)
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}

	// Encode token bytes to base64
	InitAPIKey.Token = base64.URLEncoding.EncodeToString(tokenBytes)

	c.JSON(200, InitAPIKey)
}

func GET_ServerInfo(c *gin.Context) {
	// Generate public key
	pubKey, err := GetWireguardPublicKey(ENV.WG_PRIVATE_KEY)
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}

	// Get server netmask
	mask, err := GetMask(ENV.SERVER_CIDR)
	if err != nil {
		log.Println(err)
		c.Status(500)
		return
	}

	serverInfo := types.ServerInfo{
		PublicKey:        pubKey,
		PublicEndpoint:   ENV.PUBLIC_HOST + ":" + ENV.WG_PORT,
		NameServers:      []string{strings.Split(ENV.SERVER_ADDRESS, "/")[0]},
		Netmask:          mask,
		ServerInternalIP: ENV.SERVER_ADDRESS,
	}

	c.JSON(200, serverInfo)
}
