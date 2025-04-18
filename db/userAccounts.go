package db

import (
	"github.com/wg-controller/wg-controller/types"
)

func GetAccounts() ([]types.UserAccount, error) {
	// Query the database
	query := `SELECT
		email,
		role,
		failed_attempts,
		last_active_unixmillis
		FROM user_accounts`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}

	// Loop through the rows
	var accounts []types.UserAccount
	for rows.Next() {
		var account types.UserAccount
		err = rows.Scan(
			&account.Email,
			&account.Role,
			&account.FailedAttempts,
			&account.LastActiveUnixMillis,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func GetAccount(email string) (types.UserAccount, error) {
	// Query the database
	query := `SELECT
		email,
		role,
		failed_attempts,
		last_active_unixmillis
		FROM user_accounts
		WHERE email = ?`
	row := DB.QueryRow(query, email)

	// Scan the row
	var account types.UserAccount
	err := row.Scan(
		&account.Email,
		&account.Role,
		&account.FailedAttempts,
		&account.LastActiveUnixMillis,
	)
	if err != nil {
		return types.UserAccount{}, err
	}

	return account, nil
}

func InsertAccount(email string, role string, passwordHash []byte, passwordSalt []byte) error {
	query := `INSERT INTO user_accounts (
		email,
		role,
		failed_attempts,
		password_hash,
		password_salt,
		last_active_unixmillis
	) VALUES (?, ?, ?, ?, ?, ?)`

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query,
		email,
		role,
		0,
		passwordHash,
		passwordSalt,
		0,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func UpdateAccount(account types.UserAccount) error {
	query := `UPDATE user_accounts SET
		role = ?,
		failed_attempts = ?,
		last_active_unixmillis = ?
		WHERE email = ?`

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query,
		account.Role,
		account.FailedAttempts,
		account.LastActiveUnixMillis,
		account.Email,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Delete an account from the database
func DeleteAccount(email string) error {
	query := `DELETE FROM user_accounts WHERE email = ?`

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, email)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Delete all admin accounts from the database
func DeleteAdminAccounts() error {
	query := `DELETE FROM user_accounts WHERE role = 'admin'`

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetAccountPasswordHash(email string) (hash []byte, salt []byte, err error) {
	query := `SELECT
		password_hash,
		password_salt
		FROM user_accounts
		WHERE email = ?`
	row := DB.QueryRow(query, email)

	// Scan the row
	err = row.Scan(&hash, &salt)
	return
}

func UpdateAccountPasswordHash(email string, hash []byte, salt []byte) error {
	query := `UPDATE user_accounts SET
		password_hash = ?,
		password_salt = ?
		WHERE email = ?`

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, hash, salt, email)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func IncrementAccountFailedAttempts(email string) error {
	query := `UPDATE user_accounts SET
		failed_attempts = failed_attempts + 1
		WHERE email = ?`

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, email)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func UpdateAccountLastActive(email string, unixMillis int64) error {
	query := `UPDATE user_accounts SET
		last_active_unixmillis = ?
		WHERE email = ?`

	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, unixMillis, email)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
