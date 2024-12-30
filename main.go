package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lampy255/net-tbm/db"
	"github.com/lampy255/net-tbm/types"
)

// Version
var IMAGE_TAG string

// Global Vars
var ENV types.Env

func main() {
	// Check for command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "wg-key-pair":
			privateKey, publicKey, err := NewWireguardKeyPair()
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.WriteString("WG_PRIVATE_KEY: " + string(privateKey[:]) + "\n")
			os.Stdout.WriteString("WG_PUBLIC_KEY: " + string(publicKey[:]) + "\n")
			os.Exit(0)
		case "aes-key":
			key, err := NewAESKey()
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.WriteString("DB_AES_KEY: " + key)
			os.Exit(0)
		default:
			fmt.Println("Available commands:")
			fmt.Println("  wg-key-pair:", "Generate a new Wireguard key pair")
			fmt.Println("  aes-key:", "Generate a new AES key")
			os.Exit(0)
		}
	}

	// Print version
	log.Println("Starting net-tbm:" + IMAGE_TAG)

	// Load environment variables
	LoadEnvVars()

	// Initialize the database
	db.InitDB([]byte(ENV.DB_AES_KEY))

	// Start wireguard-go
	StartWireguard()
	defer StopWireguard()

	// Init long polling
	InitLongPoll()

	// Start the API
	StartAPI()
}

func LoadEnvVars() {
	// Attempt to load .env file into environment
	godotenv.Load()

	ENV.PUBLIC_HOST = os.Getenv("PUBLIC_HOST")
	if ENV.PUBLIC_HOST == "" {
		log.Fatal("PUBLIC_HOST env variable is required")
	}

	ENV.ADMIN_EMAIL = os.Getenv("ADMIN_EMAIL")
	if ENV.ADMIN_EMAIL == "" {
		log.Fatal("ADMIN_EMAIL env variable is required")
	}

	ENV.ADMIN_PASS = os.Getenv("ADMIN_PASS")
	if ENV.ADMIN_PASS == "" {
		log.Fatal("ADMIN_PASS env variable is required")
	}

	ENV.WG_PRIVATE_KEY = os.Getenv("WG_PRIVATE_KEY")
	if ENV.WG_PRIVATE_KEY == "" {
		log.Fatal("WG_PRIVATE_KEY env variable is required. Use `net-tbm wg-key-pair` to generate one")
	} else {
		// Decode Base64
		_, err := base64.StdEncoding.DecodeString(ENV.WG_PRIVATE_KEY)
		if err != nil {
			log.Fatal("Invalid WG_PRIVATE_KEY (unable to decode base64)")
		}
	}

	ENV.WG_PUBLIC_KEY = os.Getenv("WG_PUBLIC_KEY")
	if ENV.WG_PUBLIC_KEY == "" {
		log.Fatal("WG_PUBLIC_KEY env variable is required. Use `net-tbm wg-key-pair` to generate one")
	} else {
		// Decode Base64
		_, err := base64.StdEncoding.DecodeString(ENV.WG_PRIVATE_KEY)
		if err != nil {
			log.Fatal("Invalid PRIVATE_KEY (unable to decode base64)")
		}
	}

	ENV.DB_AES_KEY = os.Getenv("DB_AES_KEY")
	if ENV.DB_AES_KEY == "" {
		log.Fatal("DB_AES_KEY env variable is required. Use `net-tbm aes-key` to generate one")
	} else {
		// Decode Base64
		bytes, err := base64.StdEncoding.DecodeString(ENV.DB_AES_KEY)
		if err != nil {
			log.Fatal("Invalid DB_AES_KEY (unable to decode base64)")
		}

		// Check if key is 32 bytes
		if len(bytes) != 32 {
			log.Fatal("Invalid DB_AES_KEY (must be 32 bytes)")
		}
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

	ENV.API_PORT = os.Getenv("API_PORT")
	if ENV.API_PORT == "" {
		log.Println("API_PORT is not set. Defaulting to 8081")
		ENV.API_PORT = "8081"
	}
}

// Creates the admin account as specified in the environment variables
// Deletes all other admin accounts
func InitAdminAccount() {
	err := db.DeleteAdminAccounts()
	if err != nil {
		log.Fatal(err)
	}

	adminAccount := types.UserAccount{
		Email: ENV.ADMIN_EMAIL,
		Role:  "admin",
	}

	salt, err := NewSalt()
	if err != nil {
		log.Fatal(err)
	}

	hash, err := HashPassword(ENV.ADMIN_PASS, salt)
	if err != nil {
		log.Fatal(err)
	}

	err = db.InsertAccount(adminAccount, hash, salt)
	if err != nil {
		log.Fatal(err)
	}
}
