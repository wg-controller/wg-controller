package db

import (
	"log"
	"time"
)

func GetSession(hash []byte) (expiresUnixMillis int64, email string, err error) {
	// Query the database
	query := `SELECT
		expires_unixmillis,
		user_email
		FROM sessions
		WHERE hash = ?`
	row := DB.QueryRow(query, hash)

	// Scan the row
	err = row.Scan(&expiresUnixMillis, &email)
	if err != nil {
		return 0, "", err
	}

	return expiresUnixMillis, email, nil
}

func CreateSession(hash []byte, userEmail string, expiresUnixMillis int64) error {
	// Insert the session into the database
	query := `INSERT INTO sessions
		(hash, user_email, expires_unixmillis)
		VALUES (?, ?, ?)`

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, hash, userEmail, expiresUnixMillis)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteUserSessions(userEmail string) error {
	// Delete the user's sessions from the database
	query := `DELETE FROM sessions WHERE user_email = ?`

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, userEmail)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteSession(hash []byte) error {
	// Delete the session from the database
	query := `DELETE FROM sessions WHERE hash = ?`

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, hash)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GarbageCollectSessions() {
	// Query the database
	query := `DELETE FROM sessions WHERE expires_unixmillis < ?`

	tx, err := DB.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	_, err = tx.Exec(query, time.Now().UnixMilli())
	if err != nil {
		log.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}
}

func SessionsGarbageCollector() {
	for {
		GarbageCollectSessions()
		time.Sleep(1 * time.Hour)
	}
}
