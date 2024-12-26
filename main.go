package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lampy255/net-tbm/types"
)

// Version
var IMAGE_TAG string

// Global Vars
var ENV types.Env

func main() {
	// Print version
	log.Println("Starting net-tbm:" + IMAGE_TAG)

	// Load environment variables
	LoadEnvVars()

	// Start wireguard-go
	StartWireguard()
	defer StopWireguard()

	go StartWebInterface()
}

func LoadEnvVars() {
	// Attempt to load .env file into environment
	godotenv.Load()

	ENV.PUBLIC_HOST = os.Getenv("PUBLIC_HOST")
	if ENV.PUBLIC_HOST == "" {
		log.Fatal("PUBLIC_HOST is required")
	}

	ENV.ADMIN_EMAIL = os.Getenv("ADMIN_EMAIL")
	if ENV.ADMIN_EMAIL == "" {
		log.Fatal("ADMIN_EMAIL is required")
	}

	ENV.ADMIN_PASS = os.Getenv("ADMIN_PASS")
	if ENV.ADMIN_PASS == "" {
		log.Fatal("ADMIN_PASS is required")
	}

	ENV.PRIVATE_KEY = os.Getenv("PRIVATE_KEY")
	if ENV.PRIVATE_KEY == "" {
		log.Fatal("PRIVATE_KEY is required")
	}

	ENV.PUBLIC_KEY = os.Getenv("PUBLIC_KEY")
	if ENV.PUBLIC_KEY == "" {
		log.Fatal("PUBLIC_KEY is required")
	}

	ENV.SERVER_CIDR = os.Getenv("SERVER_CIDR")
	if ENV.SERVER_CIDR == "" {
		log.Println("SERVER_CIDR is not set. Defaulting to 172.16.0.0/16")
		ENV.SERVER_CIDR = "172.16.0.0/16"
	} else {
		ip, _, err := net.ParseCIDR(ENV.SERVER_CIDR)
		if err != nil {
			log.Fatal("Invalid SERVER_CIDR")
		}
		if !ip.IsPrivate() {
			log.Fatal("SERVER_CIDR must be a private IP range")
		}
	}

	ENV.NAME_SERVERS = strings.Split(os.Getenv("NAME_SERVERS"), ",")
	if len(ENV.NAME_SERVERS) == 0 {
		log.Println("NAME_SERVERS is not set. Defaulting to 8.8.8.8")
		ENV.NAME_SERVERS = []string{"8.8.8.8"}
	}

	ENV.INTERFACE_NAME = os.Getenv("INTERFACE_NAME")
	if ENV.INTERFACE_NAME == "" {
		ENV.INTERFACE_NAME = "wg0"
	}

	ENV.WG_PORT = os.Getenv("WG_PORT")
	if ENV.WG_PORT == "" {
		log.Println("WG_PORT is not set. Defaulting to 51820")
		ENV.WG_PORT = "51820"
	}

	ENV.UI_PORT = os.Getenv("UI_PORT")
	if ENV.UI_PORT == "" {
		log.Println("UI_PORT is not set. Defaulting to 8080")
		ENV.UI_PORT = "8080"
	}

	ENV.API_PORT = os.Getenv("API_PORT")
	if ENV.API_PORT == "" {
		log.Println("API_PORT is not set. Defaulting to 8081")
		ENV.API_PORT = "8081"
	}
}
