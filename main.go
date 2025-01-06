package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/lampy255/wg-controller/db"
)

// Version
var IMAGE_TAG string

// Global Vars
type Env struct {
	PUBLIC_HOST      string   // Public host for web interface
	ADMIN_EMAIL      string   // Admin email
	ADMIN_PASS       string   // Admin password
	WG_PRIVATE_KEY   string   // Private key for wireguard
	DB_AES_KEY       []byte   // Base64 encoded 32 Byte AES key for encrypting private keys
	SERVER_CIDR      string   // CIDR Network for tunnel addresses (optional)
	NAME_SERVERS     []string // List of public DNS servers to use (optional)
	EGRESS_INTERFACE string   // Server egress interface to masquerade traffic (optional)
	WG_INTERFACE     string   // Wireguard interface name (optional)
	WG_PORT          string   // Port for wireguard to listen on (optional)
	API_PORT         string   // Port for API to listen on (optional)
}

var ENV Env

func main() {
	// Check for command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "generate-wg-key":
			privateKey, err := NewWireguardPrivateKey()
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.WriteString(privateKey + "\n")
			os.Exit(0)
		case "generate-db-key":
			key, err := GenerateRandomString(32)
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.WriteString(key + "\n")
			os.Exit(0)
		default:
			fmt.Println("Available commands:")
			fmt.Println("  generate-wg-key:", "Generate a new Wireguard private key")
			fmt.Println("  generate-db-key:", "Generate a new AES key")
			os.Exit(0)
		}
	}

	// Print version
	log.Println("Starting wg-controller:" + IMAGE_TAG)

	// Load environment variables
	LoadEnvVars()

	// Start wireguard
	StartWireguard()
	defer StopWireguard()

	// Initialize the database
	db.InitDB([]byte(ENV.DB_AES_KEY))

	// Initialize the admin account
	InitAdminAccount()

	// Sleep to allow wireguard to start
	time.Sleep(1 * time.Second)

	// Sync wireguard configuration
	err := SyncWireguardConfiguration()
	if err != nil {
		log.Fatal("Error syncing wireguard configuration:", err)
	}

	// Init the wireguard kernel interface
	SetWireguardInterface()

	// Init networking
	InitNetworking()

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
		log.Fatal("WG_PRIVATE_KEY env variable is required. Use `wg-controller generate-wg-key` to generate one")
	}

	DB_AES_KEY := os.Getenv("DB_AES_KEY")
	if DB_AES_KEY == "" {
		log.Fatal("DB_AES_KEY env variable is required. Use `wg-controller generate-db-key` to generate one")
	} else {
		// Decode Base64
		bytes, err := base64.StdEncoding.DecodeString(DB_AES_KEY)
		if err != nil {
			log.Fatal("Invalid DB_AES_KEY (unable to decode base64)")
		}

		// Check if key is 32 bytes
		if len(bytes) != 32 {
			log.Fatal("Invalid DB_AES_KEY (must be 32 bytes)")
		}

		// Set the key
		ENV.DB_AES_KEY = bytes
	}

	ENV.SERVER_CIDR = os.Getenv("SERVER_CIDR")
	if ENV.SERVER_CIDR == "" {
		log.Println("SERVER_CIDR is not set. Defaulting to 172.16.0.0/24")
		ENV.SERVER_CIDR = "172.16.0.0/24"
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
	} else if len(ENV.NAME_SERVERS) == 1 {
		if ENV.NAME_SERVERS[0] == "" {
			log.Println("NAME_SERVERS is not set. Defaulting to 8.8.8.8")
			ENV.NAME_SERVERS = []string{"8.8.8.8"}
		}
	}

	ENV.EGRESS_INTERFACE = os.Getenv("EGRESS_INTERFACE")
	if ENV.EGRESS_INTERFACE == "" {
		log.Println("EGRESS_INTERFACE is not set. Defaulting to eth0")
		ENV.EGRESS_INTERFACE = "eth0"
	}

	ENV.WG_INTERFACE = os.Getenv("WG_INTERFACE")
	if ENV.WG_INTERFACE == "" {
		log.Println("WG_INTERFACE is not set. Defaulting to wg0")
		ENV.WG_INTERFACE = "wg0"
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

	salt, err := NewSalt()
	if err != nil {
		log.Fatal(err)
	}

	hash, err := GenerateDeterministicHash([]byte(ENV.ADMIN_PASS), salt)
	if err != nil {
		log.Fatal(err)
	}

	err = db.InsertAccount(ENV.ADMIN_EMAIL, "admin", hash, salt)
	if err != nil {
		log.Fatal(err)
	}
}
