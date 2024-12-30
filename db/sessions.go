package db

import (
	"log"
	"time"
)

func GetSession(hash string) (expiresUnixMillis int64, err error) {
	// Query the database
	query := `SELECT
		expires_unix_millis
		FROM sessions
		WHERE hash = ?`
	row := DB.QueryRow(query, hash)

	// Scan the row
	err = row.Scan(&expiresUnixMillis)
	if err != nil {
		return 0, err
	}

	return expiresUnixMillis, nil
}

func CreateSession(hash string, userEmail string, expiresUnixMillis int64) error {
	// Insert the session into the database
	query := `INSERT INTO sessions
		(hash, user_email, expires_unix_millis)
		VALUES (?, ?, ?)`
	_, err := DB.Exec(query, hash, userEmail, expiresUnixMillis)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserSessions(userEmail string) error {
	// Delete the user's sessions from the database
	query := `DELETE FROM sessions WHERE user_email = ?`
	_, err := DB.Exec(query, userEmail)
	if err != nil {
		return err
	}

	return nil
}

func GarbageCollectSessions() {
	// Query the database
	query := `DELETE FROM sessions WHERE expires_unix_millis < ?`
	_, err := DB.Exec(query, time.Now().UnixMilli())
	if err != nil {
		log.Println(err)
	}
}

func SessionsGarbageCollector() {
	for {
		GarbageCollectSessions()
		time.Sleep(1 * time.Hour)
	}
}
