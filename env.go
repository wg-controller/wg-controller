package main

import (
	"encoding/base64"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PUBLIC_HOST      string // Public host for web interface
	ADMIN_EMAIL      string // Admin email
	ADMIN_PASS       string // Admin password
	WG_PRIVATE_KEY   string // Private key for wireguard
	DB_AES_KEY       []byte // Base64 encoded 32 Byte AES key for encrypting private keys
	SERVER_CIDR      string // CIDR Network for tunnel addresses (optional)
	SERVER_ADDRESS   string // Internal IP address of the server
	EGRESS_INTERFACE string // Server egress interface to masquerade traffic (optional)
	WG_INTERFACE     string // Wireguard interface name (optional)
	WG_PORT          string // Port for wireguard to listen on (optional)
	API_PORT         string // Port for API to listen on (optional)
	SERVER_HOSTNAME  string // Internal hostname of the server (optional)
	UPSTREAM_DNS     string // Upstream DNS server (optional)
	SLACK_WEBHOOK    string // Slack webhook URL (optional)
	PING_MONITORING  bool   // Enable ping monitoring (optional)
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
		log.Println("SERVER_CIDR is not set. Defaulting to 172.19.0.0/24")
		ENV.SERVER_CIDR = "172.19.0.0/24"
	} else {
		ip, _, err := net.ParseCIDR(ENV.SERVER_CIDR)
		if err != nil {
			log.Fatal("Invalid SERVER_CIDR")
		}
		if !ip.IsPrivate() {
			log.Fatal("SERVER_CIDR must be a private IP range")
		}
	}

	ENV.SERVER_ADDRESS = os.Getenv("SERVER_ADDRESS")
	if ENV.SERVER_ADDRESS == "" {
		log.Println("SERVER_ADDRESS is not set. Defaulting to CIDR max address")
		addr, mask, err := HighestIP(ENV.SERVER_CIDR)
		ENV.SERVER_ADDRESS = addr + mask
		if err != nil {
			log.Fatal(err)
		}
	} else {
		ip := net.ParseIP(ENV.SERVER_ADDRESS)
		if ip == nil {
			log.Fatal("Invalid SERVER_ADDRESS")
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

	ENV.SERVER_HOSTNAME = os.Getenv("SERVER_HOSTNAME")
	if ENV.SERVER_HOSTNAME == "" {
		log.Println("SERVER_HOSTNAME is not set. Defaulting to wg-controller")
		ENV.SERVER_HOSTNAME = "wg-controller"
	}

	ENV.UPSTREAM_DNS = os.Getenv("UPSTREAM_DNS")
	if ENV.UPSTREAM_DNS == "" {
		log.Println("UPSTREAM_DNS is not set. Defaulting to 8.8.8.8")
		ENV.UPSTREAM_DNS = "8.8.8.8"
	}

	ENV.SLACK_WEBHOOK = os.Getenv("SLACK_WEBHOOK")

	ENV.PING_MONITORING = os.Getenv("PING_MONITORING") == "true"
	if ENV.PING_MONITORING {
		log.Println("Internal ping monitoring enabled")
	}
}
