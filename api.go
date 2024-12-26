package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func StartWebInterface() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.New()

	// Define endpoints
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Start Listening
	log.Println("Starting web interface on port " + ENV.UI_PORT)
	err := router.Run(":" + ENV.UI_PORT)
	if err != nil {
		log.Fatal("Error starting API:", err)
	}
}

func StartClientAPI() {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	router := gin.New()

	// Define endpoints
	router.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.GET("/api/v1/peers", API_Get_Peers)

	// Start Listening
	log.Println("Starting Client API on port " + ENV.API_PORT)
	err := router.Run(":" + ENV.API_PORT)
	if err != nil {
		log.Fatal("Error starting API:", err)
	}
}

func API_Get_Peers(c *gin.Context) {
	peers, err := GetPeers()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	extendedPeers, err := GetWireguardPeers(peers)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, extendedPeers)
}
