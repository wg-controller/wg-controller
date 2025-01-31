package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wg-controller/wg-controller/db"
)

// Version
var IMAGE_TAG string

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
	time.Sleep(2 * time.Second)

	// Sync wireguard configuration
	err := SyncWireguardConfiguration()
	if err != nil {
		log.Fatal("Error syncing wireguard configuration:", err)
	}

	// Init the wireguard kernel interface
	SetWireguardInterface()

	// Init networking
	InitNetworking()

	// Init DNS
	InitDNS()

	// Sync routing table
	err = SyncRoutingTable()
	if err != nil {
		log.Fatal("Error syncing routing table:", err)
	}

	// Init ping monitoring
	if ENV.PING_MONITORING {
		go InitInternalPing()
	}

	// Init long polling
	InitLongPoll()

	// Start the API
	StartAPI()
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
