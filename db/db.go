package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var AES_KEY []byte

func InitDB(EncryptionKey []byte) {
	// Store AES key
	AES_KEY = EncryptionKey

	// Does the data directory exist?
	_, err := os.Stat("data")
	if err != nil {
		// Create the data directory
		err = os.Mkdir("data", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Open the database
	path := filepath.Join("data", "tbm-server.db")
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the peers table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS peers (
		uuid TEXT PRIMARY KEY,
		hostname TEXT UNIQUE,
		enabled BOOLEAN,
		peer_type TEXT,
		updated_unixmillis INTEGER,
		private_key TEXT,
		public_key TEXT,
		pre_shared_key TEXT,
		keep_alive_unixmillis INTEGER,
		local_tun_address TEXT,
		remote_tun_address TEXT,
		remote_subnets TEXT,
		allowed_subnets TEXT,
		last_seen_unixmillis INTEGER,
		last_ip_address TEXT,
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the user_accounts table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user_accounts (
		email TEXT PRIMARY KEY,
		role TEXT,
		failed_attempts INTEGER,
		password_hash TEXT,
		password_salt TEXT,
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the sessions table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		hash TEXT PRIMARY KEY,
		expires_unixmillis INTEGER,
		user_email TEXT,
		role TEXT,
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Create the api_keys table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS api_keys (
		uuid TEXT PRIMARY KEY,
		hash TEXT,
		expires_unixmillis INTEGER,
		role TEXT,
		name TEXT,
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Update the global DB variable
	DB = db

	// Init sessions garbage collector
	go SessionsGarbageCollector()
}
